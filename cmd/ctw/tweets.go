package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	publish "github.com/0dayfall/ctw/internal/tweet/publish"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newTweetsCommand())
}

func newTweetsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tweets",
		Short: "Manage tweets",
	}

	cmd.AddCommand(newTweetsCreateCommand())
	cmd.AddCommand(newTweetsDeleteCommand())

	return cmd
}

func newTweetsCreateCommand() *cobra.Command {
	var (
		text     string
		filePath string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a tweet",
		RunE: func(cmd *cobra.Command, args []string) error {
			if text == "" && filePath == "" {
				return errors.New("provide --text or --file")
			}

			if text != "" && filePath != "" {
				return errors.New("use either --text or --file, not both")
			}

			if filePath != "" {
				contents, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("read file: %w", err)
				}
				text = strings.TrimSpace(string(contents))
			}

			if text == "" {
				return errors.New("tweet text is empty")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := publish.NewService(c)
			response, rateLimits, err := service.CreateTweet(ctx, publish.CreateTweetRequest{Text: text})
			if err != nil {
				return err
			}

			if err := printJSON(response); err != nil {
				return err
			}
			printRateLimits(rateLimits)
			return nil
		},
	}

	cmd.Flags().StringVar(&text, "text", "", "Tweet text content")
	cmd.Flags().StringVar(&filePath, "file", "", "Path to file containing tweet text")

	return cmd
}

func newTweetsDeleteCommand() *cobra.Command {
	var tweetID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a tweet",
		RunE: func(cmd *cobra.Command, args []string) error {
			if tweetID == "" {
				return errors.New("--id is required")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := publish.NewService(c)
			response, rateLimits, err := service.DeleteTweet(ctx, tweetID)
			if err != nil {
				return err
			}

			if err := printJSON(response); err != nil {
				return err
			}
			printRateLimits(rateLimits)
			return nil
		},
	}

	cmd.Flags().StringVar(&tweetID, "id", "", "ID of the tweet to delete")

	return cmd
}
