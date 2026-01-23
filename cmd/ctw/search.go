package main

import (
	"context"
	"errors"

	recentsearch "github.com/0dayfall/ctw/internal/tweet/recentsearch"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newSearchCommand())
}

func newSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for tweets",
	}

	cmd.AddCommand(newSearchRecentCommand())
	return cmd
}

func newSearchRecentCommand() *cobra.Command {
	var (
		query      string
		nextToken  string
		extraPairs []string
	)

	cmd := &cobra.Command{
		Use:   "recent",
		Short: "Call the /2/tweets/search/recent endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			if query == "" {
				return errors.New("query is required")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			params := make(map[string]string)
			if nextToken != "" {
				params["pagination_token"] = nextToken
			}

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

			service := recentsearch.NewService(c)
			response, rateLimits, err := service.SearchRecent(ctx, query, params)
			if err != nil {
				printRateLimits(rateLimits)
				return err
			}

			if err := printJSON(response); err != nil {
				return err
			}
			printRateLimits(rateLimits)
			return nil
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Query string to search for")
	cmd.Flags().StringVar(&nextToken, "next-token", "", "Pagination token to continue a previous search")
	cmd.Flags().StringArrayVar(&extraPairs, "param", nil, "Additional query parameter in key=value format")

	return cmd
}
