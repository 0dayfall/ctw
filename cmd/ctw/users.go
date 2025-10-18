package main

import (
	"context"
	"errors"

	lookupsvc "github.com/0dayfall/ctw/internal/users/lookup"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newUsersCommand())
}

func newUsersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Work with user endpoints",
	}

	cmd.AddCommand(newUsersLookupCommand())
	cmd.AddCommand(newUsersBlockCommand())
	cmd.AddCommand(newUsersUnblockCommand())
	cmd.AddCommand(newUsersFollowCommand())
	cmd.AddCommand(newUsersUnfollowCommand())

	return cmd
}

func newUsersLookupCommand() *cobra.Command {
	var (
		id         string
		username   string
		ids        []string
		usernames  []string
		extraPairs []string
	)

	cmd := &cobra.Command{
		Use:   "lookup",
		Short: "Lookup users by id or username",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			service := lookupsvc.NewService(c)

			switch {
			case len(ids) > 0:
				response, rateLimits, err := service.LookupIDs(ctx, ids, params)
				if err != nil {
					return err
				}
				if err := printJSON(response); err != nil {
					return err
				}
				printRateLimits(rateLimits)
				return nil
			case len(usernames) > 0:
				response, rateLimits, err := service.LookupUsernames(ctx, usernames, params)
				if err != nil {
					return err
				}
				if err := printJSON(response); err != nil {
					return err
				}
				printRateLimits(rateLimits)
				return nil
			case id != "":
				response, rateLimits, err := service.LookupID(ctx, id, params)
				if err != nil {
					return err
				}
				if err := printJSON(response); err != nil {
					return err
				}
				printRateLimits(rateLimits)
				return nil
			case username != "":
				response, rateLimits, err := service.LookupUsername(ctx, username, params)
				if err != nil {
					return err
				}
				if err := printJSON(response); err != nil {
					return err
				}
				printRateLimits(rateLimits)
				return nil
			default:
				return errors.New("provide --id, --username, --ids, or --usernames")
			}
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "Single user id to lookup")
	cmd.Flags().StringVar(&username, "username", "", "Single username to lookup")
	cmd.Flags().StringSliceVar(&ids, "ids", nil, "Comma-separated list of user ids to lookup")
	cmd.Flags().StringSliceVar(&usernames, "usernames", nil, "Comma-separated list of usernames to lookup")
	cmd.Flags().StringArrayVar(&extraPairs, "param", nil, "Additional query parameter in key=value format")

	return cmd
}

func newUsersBlockCommand() *cobra.Command {
	var (
		sourceID string
		targetID string
	)

	cmd := &cobra.Command{
		Use:   "block",
		Short: "Block a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if sourceID == "" || targetID == "" {
				return errors.New("both --source-id and --target-id are required")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}
			service := lookupsvc.NewService(c)

			response, rateLimits, err := service.Block(ctx, sourceID, targetID)
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

	cmd.Flags().StringVar(&sourceID, "source-id", "", "User id performing the block")
	cmd.Flags().StringVar(&targetID, "target-id", "", "User id being blocked")

	return cmd
}

func newUsersUnblockCommand() *cobra.Command {
	var (
		sourceID string
		targetID string
	)

	cmd := &cobra.Command{
		Use:   "unblock",
		Short: "Remove an existing block",
		RunE: func(cmd *cobra.Command, args []string) error {
			if sourceID == "" || targetID == "" {
				return errors.New("both --source-id and --target-id are required")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}
			service := lookupsvc.NewService(c)

			response, rateLimits, err := service.Unblock(ctx, sourceID, targetID)
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

	cmd.Flags().StringVar(&sourceID, "source-id", "", "User id lifting the block")
	cmd.Flags().StringVar(&targetID, "target-id", "", "User id being unblocked")

	return cmd
}

func newUsersFollowCommand() *cobra.Command {
	var (
		sourceID string
		targetID string
	)

	cmd := &cobra.Command{
		Use:   "follow",
		Short: "Follow a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if sourceID == "" || targetID == "" {
				return errors.New("both --source-id and --target-id are required")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}
			service := lookupsvc.NewService(c)

			response, rateLimits, err := service.Follow(ctx, sourceID, targetID)
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

	cmd.Flags().StringVar(&sourceID, "source-id", "", "User id initiating the follow")
	cmd.Flags().StringVar(&targetID, "target-id", "", "User id being followed")

	return cmd
}

func newUsersUnfollowCommand() *cobra.Command {
	var (
		sourceID string
		targetID string
	)

	cmd := &cobra.Command{
		Use:   "unfollow",
		Short: "Stop following a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if sourceID == "" || targetID == "" {
				return errors.New("both --source-id and --target-id are required")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}
			service := lookupsvc.NewService(c)

			response, rateLimits, err := service.Unfollow(ctx, sourceID, targetID)
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

	cmd.Flags().StringVar(&sourceID, "source-id", "", "User id stopping the follow")
	cmd.Flags().StringVar(&targetID, "target-id", "", "User id being unfollowed")

	return cmd
}
