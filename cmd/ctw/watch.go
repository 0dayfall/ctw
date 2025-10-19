package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	stream "github.com/0dayfall/ctw/internal/tweet/filteredstream"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newWatchCommand())
}

func newWatchCommand() *cobra.Command {
	var (
		keywords  []string
		autoSetup bool
		showUser  bool
		showMeta  bool
	)

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch tweets in real-time for keywords",
		Long: `Watch tweets matching keywords using Twitter's filtered stream.

This command sets up stream rules and monitors tweets in real-time. 
Use Ctrl+C to stop watching.

Examples:
  # Watch for specific keywords
  ctw watch --keyword "golang" --keyword "rust"

  # Auto-setup stream rules
  ctw watch --keyword "AI" --auto-setup

  # Show detailed information
  ctw watch --keyword "bitcoin" --show-user --show-meta`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(keywords) == 0 {
				return errors.New("at least one keyword is required (use --keyword)")
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Handle Ctrl+C gracefully
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-sigChan
				fmt.Fprintf(os.Stderr, "\n\n🛑 Stopping stream...\n")
				cancel()
			}()

			c, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := stream.NewService(c)

			// Auto-setup rules if requested
			if autoSetup {
				fmt.Fprintf(os.Stderr, "📝 Setting up stream rules...\n")

				// Clear existing rules first
				existing, _, listErr := service.GetRules(ctx)
				if listErr != nil {
					return fmt.Errorf("failed to get existing rules: %w", listErr)
				}

				if len(existing.Data) > 0 {
					ids := make([]string, len(existing.Data))
					for i, rule := range existing.Data {
						ids[i] = rule.ID
					}
					deleteCmd := stream.CreateDeleteIdCommand(ids)
					if _, _, delErr := service.DeleteRule(ctx, deleteCmd, false); delErr != nil {
						fmt.Fprintf(os.Stderr, "⚠️  Warning: failed to delete existing rules: %v\n", delErr)
					}
				}
				// Add new rules for each keyword
				adds := make([]stream.Add, 0, len(keywords))
				for _, keyword := range keywords {
					adds = append(adds, stream.Add{
						Value: keyword,
						Tag:   fmt.Sprintf("watch_%s", strings.ReplaceAll(keyword, " ", "_")),
					})
				}

				addCmd := stream.AddCommand{Add: adds}
				resp, _, err := service.AddRule(ctx, addCmd, false)
				if err != nil {
					return fmt.Errorf("failed to add rules: %w", err)
				}

				fmt.Fprintf(os.Stderr, "✅ Added %d rule(s)\n", len(resp.Data))
				for _, rule := range resp.Data {
					fmt.Fprintf(os.Stderr, "   - %s (ID: %s)\n", rule.Value, rule.ID)
				}
			}

			// Build stream fields
			fields := map[string]string{
				"tweet.fields": "created_at,author_id,lang,possibly_sensitive,source",
			}
			if showUser {
				fields["expansions"] = "author_id"
				fields["user.fields"] = "name,username,created_at"
			}

			fmt.Fprintf(os.Stderr, "\n🔴 Watching for keywords: %s\n", strings.Join(keywords, ", "))
			fmt.Fprintf(os.Stderr, "Press Ctrl+C to stop\n\n")

			tweetCount := 0
			startTime := time.Now()

			// Stream tweets
			err = service.StreamReader(ctx, fields, func(tweet stream.StreamTweet, includes stream.StreamIncludes) error {
				tweetCount++

				// Display tweet
				fmt.Printf("\n─────────────────────────────────────────\n")
				fmt.Printf("🐦 Tweet #%d\n", tweetCount)
				fmt.Printf("ID: %s\n", tweet.ID)
				fmt.Printf("Time: %s\n", tweet.CreatedAt.Format(time.RFC3339))

				if showUser && len(includes.Users) > 0 {
					user := includes.Users[0]
					fmt.Printf("Author: @%s (%s)\n", user.Username, user.Name)
				} else {
					fmt.Printf("Author ID: %s\n", tweet.AuthorID)
				}

				fmt.Printf("Language: %s\n", tweet.Lang)
				if tweet.PossiblySensitive {
					fmt.Printf("⚠️  Possibly Sensitive\n")
				}

				fmt.Printf("\nText:\n%s\n", tweet.Text)

				if showMeta {
					fmt.Printf("\nMetadata:\n")
					fmt.Printf("  Source: %s\n", tweet.Source)
				}

				return nil
			})

			if err != nil && err != context.Canceled {
				return fmt.Errorf("stream error: %w", err)
			}

			// Show summary
			duration := time.Since(startTime)
			fmt.Fprintf(os.Stderr, "\n\n📊 Summary:\n")
			fmt.Fprintf(os.Stderr, "   Tweets received: %d\n", tweetCount)
			fmt.Fprintf(os.Stderr, "   Duration: %s\n", duration.Round(time.Second))
			if duration.Seconds() > 0 {
				rate := float64(tweetCount) / duration.Seconds() * 60
				fmt.Fprintf(os.Stderr, "   Rate: %.1f tweets/minute\n", rate)
			}

			return nil
		},
	}

	cmd.Flags().StringArrayVar(&keywords, "keyword", nil, "Keyword to watch for (can be specified multiple times)")
	cmd.Flags().BoolVar(&autoSetup, "auto-setup", false, "Automatically set up stream rules for keywords")
	cmd.Flags().BoolVar(&showUser, "show-user", false, "Show author information")
	cmd.Flags().BoolVar(&showMeta, "show-meta", false, "Show additional metadata")

	return cmd
}
