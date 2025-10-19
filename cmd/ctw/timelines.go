package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/0dayfall/ctw/internal/tweet/timelines"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newTimelinesCommand())
}

func newTimelinesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timelines",
		Short: "Fetch user timelines",
	}

	cmd.AddCommand(newTimelinesUserCommand())
	cmd.AddCommand(newTimelinesMentionsCommand())
	cmd.AddCommand(newTimelinesHomeCommand())

	return cmd
}

func newTimelinesUserCommand() *cobra.Command {
	var (
		userID     string
		paramsFlag []string
	)

	cmd := &cobra.Command{
		Use:   "user",
		Short: "Get tweets posted by a user",
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

			service := timelines.NewService(c)
			response, rateLimits, err := service.GetUserTweets(ctx, userID, queryParams)
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

	cmd.Flags().StringVar(&userID, "user-id", "", "ID of the user")
	cmd.Flags().StringSliceVar(&paramsFlag, "param", nil, "Additional query parameters in key=value form (repeatable)")

	return cmd
}

func newTimelinesMentionsCommand() *cobra.Command {
	var (
		userID     string
		paramsFlag []string
	)

	cmd := &cobra.Command{
		Use:   "mentions",
		Short: "Get tweets that mention a user",
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

			service := timelines.NewService(c)
			response, rateLimits, err := service.GetUserMentions(ctx, userID, queryParams)
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

	cmd.Flags().StringVar(&userID, "user-id", "", "ID of the user")
	cmd.Flags().StringSliceVar(&paramsFlag, "param", nil, "Additional query parameters in key=value form (repeatable)")

	return cmd
}

func newTimelinesHomeCommand() *cobra.Command {
	var (
		userID     string
		paramsFlag []string
	)

	cmd := &cobra.Command{
		Use:   "home",
		Short: "Get reverse chronological home timeline for authenticated user",
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

			service := timelines.NewService(c)
			response, rateLimits, err := service.GetReverseChronological(ctx, userID, queryParams)
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

	cmd.Flags().StringVar(&userID, "user-id", "", "ID of the authenticated user")
	cmd.Flags().StringSliceVar(&paramsFlag, "param", nil, "Additional query parameters in key=value form (repeatable)")

	return cmd
}
