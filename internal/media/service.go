package media

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	defaultUploadBaseURL = "https://upload.twitter.com/1.1/media/"
	chunkSize            = 5 * 1024 * 1024 // 5MB chunks
)

// Service coordinates Twitter media upload operations.
type Service struct {
	bearerToken   string
	httpClient    *http.Client
	uploadBaseURL string
}

// NewService constructs a Service with the provided bearer token.
func NewService(bearerToken string) *Service {
	return &Service{
		bearerToken:   bearerToken,
		httpClient:    &http.Client{Timeout: 120 * time.Second},
		uploadBaseURL: defaultUploadBaseURL,
	}
}

// UploadFile uploads a media file using the chunked upload flow.
// For large files or videos, this uses INIT -> APPEND -> FINALIZE.
// For small images, it may use simple upload.
func (s *Service) UploadFile(ctx context.Context, filePath string, category MediaCategory) (string, error) {
	if s == nil {
		return "", fmt.Errorf("media: nil service")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("media: open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("media: stat file: %w", err)
	}

	fileSize := stat.Size()
	mediaType := detectMediaType(filePath)

	// INIT phase
	initResp, err := s.initUpload(ctx, fileSize, mediaType, category)
	if err != nil {
		return "", fmt.Errorf("media: init upload: %w", err)
	}

	// APPEND phase - upload in chunks
	err = s.appendChunks(ctx, file, initResp.MediaIDString, fileSize)
	if err != nil {
		return "", fmt.Errorf("media: append chunks: %w", err)
	}

	// FINALIZE phase
	finalizeResp, err := s.finalizeUpload(ctx, initResp.MediaIDString)
	if err != nil {
		return "", fmt.Errorf("media: finalize upload: %w", err)
	}

	// Wait for async processing if needed (videos)
	if finalizeResp.ProcessingInfo != nil {
		if err := s.waitForProcessing(ctx, initResp.MediaIDString); err != nil {
			return "", fmt.Errorf("media: wait for processing: %w", err)
		}
	}

	return finalizeResp.MediaIDString, nil
}

func (s *Service) initUpload(ctx context.Context, totalBytes int64, mediaType string, category MediaCategory) (*InitResponse, error) {
	params := url.Values{}
	params.Set("command", "INIT")
	params.Set("total_bytes", strconv.FormatInt(totalBytes, 10))
	params.Set("media_type", mediaType)
	if category != "" {
		params.Set("media_category", string(category))
	}

	req, err := s.newRequest(ctx, "POST", "upload.json", params, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("init upload failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var initResp InitResponse
	if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
		return nil, fmt.Errorf("decode init response: %w", err)
	}

	return &initResp, nil
}

func (s *Service) appendChunks(ctx context.Context, file *os.File, mediaID string, fileSize int64) error {
	segmentIndex := 0
	buffer := make([]byte, chunkSize)

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read chunk: %w", err)
		}

		chunk := buffer[:n]
		if err := s.appendChunk(ctx, mediaID, segmentIndex, chunk); err != nil {
			return fmt.Errorf("append chunk %d: %w", segmentIndex, err)
		}

		segmentIndex++
	}

	return nil
}

func (s *Service) appendChunk(ctx context.Context, mediaID string, segmentIndex int, data []byte) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	writer.WriteField("command", "APPEND")
	writer.WriteField("media_id", mediaID)
	writer.WriteField("segment_index", strconv.Itoa(segmentIndex))

	part, err := writer.CreateFormFile("media", "chunk")
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}
	_, err = part.Write(data)
	if err != nil {
		return fmt.Errorf("write chunk data: %w", err)
	}

	contentType := writer.FormDataContentType()
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.uploadBaseURL+"upload.json", &body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.bearerToken)
	req.Header.Set("Content-Type", contentType)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("append chunk failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *Service) finalizeUpload(ctx context.Context, mediaID string) (*FinalizeResponse, error) {
	params := url.Values{}
	params.Set("command", "FINALIZE")
	params.Set("media_id", mediaID)

	req, err := s.newRequest(ctx, "POST", "upload.json", params, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("finalize upload failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var finalizeResp FinalizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&finalizeResp); err != nil {
		return nil, fmt.Errorf("decode finalize response: %w", err)
	}

	return &finalizeResp, nil
}

func (s *Service) waitForProcessing(ctx context.Context, mediaID string) error {
	maxAttempts := 60
	attempt := 0

	for attempt < maxAttempts {
		status, err := s.checkStatus(ctx, mediaID)
		if err != nil {
			return err
		}

		if status.ProcessingInfo == nil {
			return nil
		}

		switch status.ProcessingInfo.State {
		case "succeeded":
			return nil
		case "failed":
			if status.ProcessingInfo.Error != nil {
				return fmt.Errorf("processing failed: %s", status.ProcessingInfo.Error.Message)
			}
			return fmt.Errorf("processing failed with no error details")
		case "in_progress", "pending":
			waitTime := time.Duration(status.ProcessingInfo.CheckAfterSecs) * time.Second
			if waitTime == 0 {
				waitTime = 5 * time.Second
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(waitTime):
				attempt++
			}
		default:
			return fmt.Errorf("unknown processing state: %s", status.ProcessingInfo.State)
		}
	}

	return fmt.Errorf("processing timeout after %d attempts", maxAttempts)
}

func (s *Service) checkStatus(ctx context.Context, mediaID string) (*StatusResponse, error) {
	params := url.Values{}
	params.Set("command", "STATUS")
	params.Set("media_id", mediaID)

	req, err := s.newRequest(ctx, "GET", "upload.json", params, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status check failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var statusResp StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("decode status response: %w", err)
	}

	return &statusResp, nil
}

func (s *Service) newRequest(ctx context.Context, method, path string, params url.Values, body io.Reader) (*http.Request, error) {
	u := s.uploadBaseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.bearerToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}

func detectMediaType(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	default:
		return "application/octet-stream"
	}
}
