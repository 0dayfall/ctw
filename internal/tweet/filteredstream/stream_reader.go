package tweet

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/0dayfall/ctw/internal/client"
)

// TweetHandler is a function that processes tweets from the stream.
// Return an error to stop the stream.
type TweetHandler func(tweet StreamTweet, includes StreamIncludes) error

// StreamReader connects to the filtered stream and processes tweets in real-time.
func (s *Service) StreamReader(ctx context.Context, fields map[string]string, handler TweetHandler) error {
	resp, err := s.client.Get(ctx, streamPath, fields)
	if err != nil {
		return err
	}
	defer client.SafeClose(resp.Body)

	if err := client.CheckResponse(resp); err != nil {
		return err
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024) // 1MB max token size

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Bytes()
		if len(line) == 0 {
			continue // Skip empty lines (keep-alive)
		}

		var envelope StreamEnvelope
		if err := json.Unmarshal(line, &envelope); err != nil {
			// Log and continue on parse errors
			fmt.Printf("Error parsing tweet: %v\n", err)
			continue
		}

		// Process each tweet in the response
		for _, tweet := range envelope.Data {
			if err := handler(tweet, envelope.Includes); err != nil {
				if err == io.EOF {
					return nil // Clean stop
				}
				return fmt.Errorf("tweet handler error: %w", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stream scanner error: %w", err)
	}

	return nil
}
