package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/0dayfall/ctw/internal/tweet/likes"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newLikesCommand())
}

func newLikesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "likes",
		Short: "Manage tweet likes",
	}

	cmd.AddCommand(newLikesAddCommand())
	cmd.AddCommand(newLikesRemoveCommand())
	cmd.AddCommand(newLikesListCommand())

	return cmd
}

func newLikesAddCommand() *cobra.Command {
	var (
		userID  string
		tweetID string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Like a tweet",
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

			service := likes.NewService(c)
			response, rateLimits, err := service.LikeTweet(ctx, userID, tweetID)
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

	cmd.Flags().StringVar(&userID, "user-id", "", "ID of the user liking the tweet")
	cmd.Flags().StringVar(&tweetID, "tweet-id", "", "ID of the tweet to like")

	return cmd
}

func newLikesRemoveCommand() *cobra.Command {
	var (
		userID  string
		tweetID string
	)

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Unlike a tweet",
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

			service := likes.NewService(c)
			response, rateLimits, err := service.UnlikeTweet(ctx, userID, tweetID)
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

	cmd.Flags().StringVar(&userID, "user-id", "", "ID of the user unliking the tweet")
	cmd.Flags().StringVar(&tweetID, "tweet-id", "", "ID of the tweet to unlike")

	return cmd
}

func newLikesListCommand() *cobra.Command {
	var (
		userID     string
		paramsFlag []string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List liked tweets for a user",
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

			service := likes.NewService(c)
			response, rateLimits, err := service.ListLikedTweets(ctx, userID, queryParams)
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

	cmd.Flags().StringVar(&userID, "user-id", "", "ID of the user whose likes to list")
	cmd.Flags().StringSliceVar(&paramsFlag, "param", nil, "Additional query parameters in key=value form (repeatable)")

	return cmd
}
