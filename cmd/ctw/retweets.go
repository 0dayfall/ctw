package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/0dayfall/ctw/internal/tweet/retweets"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newRetweetsCommand())
}

func newRetweetsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retweets",
		Short: "Manage retweets",
	}

	cmd.AddCommand(newRetweetsAddCommand())
	cmd.AddCommand(newRetweetsRemoveCommand())
	cmd.AddCommand(newRetweetsListCommand())

	return cmd
}

func newRetweetsAddCommand() *cobra.Command {
	var (
		userID  string
		tweetID string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Retweet a tweet",
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

			service := retweets.NewService(c)
			response, rateLimits, err := service.Retweet(ctx, userID, tweetID)
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

	cmd.Flags().StringVar(&userID, "user-id", "", "ID of the user retweeting")
	cmd.Flags().StringVar(&tweetID, "tweet-id", "", "ID of the tweet to retweet")

	return cmd
}

func newRetweetsRemoveCommand() *cobra.Command {
	var (
		userID  string
		tweetID string
	)

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a retweet",
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

			service := retweets.NewService(c)
			response, rateLimits, err := service.Unretweet(ctx, userID, tweetID)
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

	cmd.Flags().StringVar(&userID, "user-id", "", "ID of the user removing the retweet")
	cmd.Flags().StringVar(&tweetID, "tweet-id", "", "ID of the tweet to unretweet")

	return cmd
}

func newRetweetsListCommand() *cobra.Command {
	var (
		tweetID    string
		paramsFlag []string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users who retweeted a tweet",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(tweetID) == "" {
				return errors.New("--tweet-id is required")
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

			service := retweets.NewService(c)
			response, rateLimits, err := service.ListRetweeters(ctx, tweetID, queryParams)
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

	cmd.Flags().StringVar(&tweetID, "tweet-id", "", "ID of the tweet to list retweeters for")
	cmd.Flags().StringSliceVar(&paramsFlag, "param", nil, "Additional query parameters in key=value form (repeatable)")

	return cmd
}
