package main

import (
	"context"
	"errors"

	recentcount "github.com/0dayfall/ctw/internal/tweet/recentcount"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCountsCommand())
}

func newCountsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "counts",
		Short: "Retrieve tweet counts",
	}

	cmd.AddCommand(newCountsRecentCommand())
	cmd.AddCommand(newCountsAllCommand())
	return cmd
}

func newCountsRecentCommand() *cobra.Command {
	var (
		query       string
		granularity string
		extraPairs  []string
	)

	cmd := &cobra.Command{
		Use:   "recent",
		Short: "Call the /2/tweets/counts/recent endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			if query == "" {
				return errors.New("query is required")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			params := map[string]string{}
			if len(extraPairs) > 0 {
				extras, err := parseKeyValuePairs(extraPairs)
				if err != nil {
					return err
				}
				for k, v := range extras {
					params[k] = v
				}
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := recentcount.NewService(c)
			response, rateLimits, err := service.GetRecentCount(ctx, query, granularity, params)
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

	cmd.Flags().StringVar(&query, "query", "", "Query to aggregate counts for")
	cmd.Flags().StringVar(&granularity, "granularity", "day", "Granularity (minute|hour|day)")
	cmd.Flags().StringArrayVar(&extraPairs, "param", nil, "Additional query parameter in key=value format")

	return cmd
}

func newCountsAllCommand() *cobra.Command {
	var (
		query       string
		granularity string
		extraPairs  []string
	)

	cmd := &cobra.Command{
		Use:   "all",
		Short: "Call the /2/tweets/counts/all endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			if query == "" {
				return errors.New("query is required")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			params := map[string]string{}
			if len(extraPairs) > 0 {
				extras, err := parseKeyValuePairs(extraPairs)
				if err != nil {
					return err
				}
				for k, v := range extras {
					params[k] = v
				}
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := recentcount.NewService(c)
			response, rateLimits, err := service.GetAllCount(ctx, query, granularity, params)
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

	cmd.Flags().StringVar(&query, "query", "", "Query to aggregate counts for")
	cmd.Flags().StringVar(&granularity, "granularity", "day", "Granularity (minute|hour|day)")
	cmd.Flags().StringArrayVar(&extraPairs, "param", nil, "Additional query parameter in key=value format")

	return cmd
}
