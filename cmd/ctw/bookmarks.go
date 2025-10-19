package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/0dayfall/ctw/internal/tweet/bookmarks"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newBookmarksCommand())
}

func newBookmarksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bookmarks",
		Short: "Manage bookmarks",
	}

	cmd.AddCommand(newBookmarksAddCommand())
	cmd.AddCommand(newBookmarksRemoveCommand())
	cmd.AddCommand(newBookmarksListCommand())

	return cmd
}

func newBookmarksAddCommand() *cobra.Command {
	var (
		userID  string
		tweetID string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Bookmark a tweet",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(userID) == "" {
				return errors.New("--user-id is required")
			}
			if strings.TrimSpace(tweetID) == "" {
				return errors.New("--tweet-id is required")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := bookmarks.NewService(c)
			response, rateLimits, err := service.Add(ctx, userID, tweetID)
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

	cmd.Flags().StringVar(&userID, "user-id", "", "ID of the user bookmarking")
	cmd.Flags().StringVar(&tweetID, "tweet-id", "", "ID of the tweet to bookmark")

	return cmd
}

func newBookmarksRemoveCommand() *cobra.Command {
	var (
		userID  string
		tweetID string
	)

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a bookmark",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(userID) == "" {
				return errors.New("--user-id is required")
			}
			if strings.TrimSpace(tweetID) == "" {
				return errors.New("--tweet-id is required")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := bookmarks.NewService(c)
			response, rateLimits, err := service.Remove(ctx, userID, tweetID)
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

	cmd.Flags().StringVar(&userID, "user-id", "", "ID of the user removing the bookmark")
	cmd.Flags().StringVar(&tweetID, "tweet-id", "", "ID of the tweet to unbookmark")

	return cmd
}

func newBookmarksListCommand() *cobra.Command {
	var (
		userID     string
		paramsFlag []string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List bookmarked tweets for a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(userID) == "" {
				return errors.New("--user-id is required")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			queryParams, err := parseKeyValuePairs(paramsFlag)
			if err != nil {
				return fmt.Errorf("parse params: %w", err)
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := bookmarks.NewService(c)
			response, rateLimits, err := service.List(ctx, userID, queryParams)
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

	cmd.Flags().StringVar(&userID, "user-id", "", "ID of the user whose bookmarks to list")
	cmd.Flags().StringSliceVar(&paramsFlag, "param", nil, "Additional query parameters in key=value form (repeatable)")

	return cmd
}
