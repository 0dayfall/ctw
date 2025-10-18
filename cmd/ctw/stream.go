package main

import (
	"context"
	"errors"

	stream "github.com/0dayfall/ctw/internal/tweet/filteredstream"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newStreamCommand())
}

func newStreamCommand() *cobra.Command {
	var fieldPairs []string

	cmd := &cobra.Command{
		Use:   "stream",
		Short: "Call the filtered stream endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			fields, err := parseKeyValuePairs(fieldPairs)
			if err != nil {
				return err
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := stream.NewService(c)
			response, rateLimits, err := service.Stream(ctx, fields)
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

	cmd.Flags().StringArrayVar(&fieldPairs, "field", nil, "Query parameter to include in the request (key=value)")
	cmd.AddCommand(newStreamRulesCommand())

	return cmd
}

func newStreamRulesCommand() *cobra.Command {
	rulesCmd := &cobra.Command{
		Use:   "rules",
		Short: "Manage filtered stream rules",
	}

	rulesCmd.AddCommand(newStreamRulesAddCommand())
	rulesCmd.AddCommand(newStreamRulesListCommand())

	return rulesCmd
}

func newStreamRulesAddCommand() *cobra.Command {
	var (
		value string
		tag   string
		dry   bool
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a filtered stream rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if value == "" {
				return errors.New("rule value is required")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := stream.NewService(c)
			payload := stream.AddCommand{Add: []stream.Add{{Value: value}}}
			if tag != "" {
				payload.Add[0].Tag = tag
			}

			response, rateLimits, err := service.AddRule(ctx, payload, dry)
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

	cmd.Flags().StringVar(&value, "value", "", "Rule value to add (e.g. 'cats has:images')")
	cmd.Flags().StringVar(&tag, "tag", "", "Optional rule tag")
	cmd.Flags().BoolVar(&dry, "dry-run", false, "Send the request as a dry run to validate the rule")

	return cmd
}

func newStreamRulesListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List filtered stream rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := stream.NewService(c)
			response, rateLimits, err := service.GetRules(ctx)
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

	return cmd
}
