package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/0dayfall/ctw/internal/media"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newMediaCommand())
}

func newMediaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "media",
		Short: "Upload media files",
	}

	cmd.AddCommand(newMediaUploadCommand())

	return cmd
}

func newMediaUploadCommand() *cobra.Command {
	var (
		filePath string
		category string
	)

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload an image or video",
		Long: `Upload media to Twitter using the chunked upload API.
Supports images (JPEG, PNG, GIF, WebP) and videos (MP4, MOV).
Videos are processed asynchronously and this command will wait for completion.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(filePath) == "" {
				return errors.New("--file is required")
			}

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("file does not exist: %s", filePath)
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			// Get bearer token from flags or environment
			token := bearerTokenFlag
			if token == "" {
				token = os.Getenv("BEARER_TOKEN")
			}
			if token == "" {
				return errors.New("bearer token required (--bearer-token or BEARER_TOKEN environment variable)")
			}

			var mediaCategory media.MediaCategory
			if category != "" {
				mediaCategory = media.MediaCategory(category)
			}

			service := media.NewService(token)
			mediaID, err := service.UploadFile(ctx, filePath, mediaCategory)
			if err != nil {
				return fmt.Errorf("upload failed: %w", err)
			}

			result := map[string]string{
				"media_id_string": mediaID,
				"status":          "uploaded",
			}

			if err := printJSON(result); err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "Media uploaded successfully. Use media_id_string in tweets/DMs.\n")
			return nil
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "Path to the media file to upload (required)")
	cmd.Flags().StringVar(&category, "category", "", "Media category (tweet_image, tweet_video, dm_image, etc.)")

	return cmd
}
