// Package media provides helpers for uploading media to Twitter.
package media

// MediaCategory represents the type of media being uploaded.
type MediaCategory string

const (
	CategoryTweetImage MediaCategory = "tweet_image"
	CategoryTweetVideo MediaCategory = "tweet_video"
	CategoryTweetGif   MediaCategory = "tweet_gif"
	CategoryDMImage    MediaCategory = "dm_image"
	CategoryDMVideo    MediaCategory = "dm_video"
	CategoryDMGif      MediaCategory = "dm_gif"
)

// InitResponse captures the response from the INIT phase.
type InitResponse struct {
	MediaID          int64  `json:"media_id"`
	MediaIDString    string `json:"media_id_string"`
	ExpiresAfterSecs int    `json:"expires_after_secs,omitempty"`
}

// FinalizeResponse captures the response from the FINALIZE phase.
type FinalizeResponse struct {
	MediaID          int64           `json:"media_id"`
	MediaIDString    string          `json:"media_id_string"`
	Size             int64           `json:"size"`
	ExpiresAfterSecs int             `json:"expires_after_secs,omitempty"`
	ProcessingInfo   *ProcessingInfo `json:"processing_info,omitempty"`
}

// ProcessingInfo contains async processing status for videos.
type ProcessingInfo struct {
	State           string           `json:"state"`
	CheckAfterSecs  int              `json:"check_after_secs,omitempty"`
	ProgressPercent int              `json:"progress_percent,omitempty"`
	Error           *ProcessingError `json:"error,omitempty"`
}

// ProcessingError captures async processing errors.
type ProcessingError struct {
	Code    int    `json:"code"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

// StatusResponse captures the response from the STATUS endpoint.
type StatusResponse struct {
	MediaID        int64           `json:"media_id"`
	MediaIDString  string          `json:"media_id_string"`
	ProcessingInfo *ProcessingInfo `json:"processing_info,omitempty"`
}
