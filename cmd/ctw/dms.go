package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/0dayfall/ctw/internal/client"
	"github.com/0dayfall/ctw/internal/dm"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newDMsCommand())
}

func newDMsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dms",
		Short: "Manage direct messages",
	}

	cmd.AddCommand(newDMsSendCommand())
	cmd.AddCommand(newDMsListCommand())
	cmd.AddCommand(newDMsDeleteCommand())

	return cmd
}

func newDMsSendCommand() *cobra.Command {
	var (
		text           string
		filePath       string
		participantID  string
		conversationID string
	)

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send a direct message",
		RunE: func(cmd *cobra.Command, args []string) error {
			if participantID == "" && conversationID == "" {
				return errors.New("provide --user-id or --conversation-id")
			}
			if participantID != "" && conversationID != "" {
				return errors.New("use either --user-id or --conversation-id, not both")
			}
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

			if strings.TrimSpace(text) == "" {
				return errors.New("message text is empty")
			}

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := dm.NewService(c)

			req := dm.SendDMRequest{Text: text}
			var (
				response   dm.SendDMResponse
				rateLimits client.RateLimitSnapshot
			)

			if participantID != "" {
				response, rateLimits, err = service.SendToUser(ctx, participantID, req)
			} else {
				response, rateLimits, err = service.SendToConversation(ctx, conversationID, req)
			}
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

	cmd.Flags().StringVar(&participantID, "user-id", "", "Participant user ID for 1:1 messages")
	cmd.Flags().StringVar(&conversationID, "conversation-id", "", "Existing conversation ID")
	cmd.Flags().StringVar(&text, "text", "", "Direct message text")
	cmd.Flags().StringVar(&filePath, "file", "", "Path to file containing DM text")

	return cmd
}

func newDMsListCommand() *cobra.Command {
	var paramsFlag []string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List direct message events",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			queryParams, err := parseKeyValuePairs(paramsFlag)
			if err != nil {
				return err
			}

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := dm.NewService(c)
			response, rateLimits, err := service.ListEvents(ctx, queryParams)
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

	cmd.Flags().StringSliceVar(&paramsFlag, "param", nil, "Additional query parameters in key=value form (repeatable)")

	return cmd
}

func newDMsDeleteCommand() *cobra.Command {
	var eventID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a direct message event",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(eventID) == "" {
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

			service := dm.NewService(c)
			rateLimits, err := service.DeleteEvent(ctx, eventID)
			if err != nil {
				return err
			}

			printRateLimits(rateLimits)
			return nil
		},
	}

	cmd.Flags().StringVar(&eventID, "id", "", "DM event ID to delete")

	return cmd
}
