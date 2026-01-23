package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
		keywords   []string
		autoSetup  bool
		showUser   bool
		showMeta   bool
		jsonOutput bool
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
			reconnects := 0
			lastDisconnect := "none"
			lastRuleSet := "existing rules"
			if autoSetup {
				lastRuleSet = strings.Join(keywords, ", ")
			}

			backoff := 2 * time.Second
			maxBackoff := resolvedSettings.StreamBackoffMax
			if maxBackoff <= 0 {
				maxBackoff = 2 * time.Minute
			}

			// Stream tweets with reconnect + backoff
			for {
				err = service.StreamReader(ctx, fields, func(tweet stream.StreamTweet, includes stream.StreamIncludes) error {
					tweetCount++

					if jsonOutput {
						if cmd.Flags().Changed("pretty") {
							type prettyTweet struct {
								ID                string    `json:"id"`
								Text              string    `json:"text"`
								AuthorID          string    `json:"author_id,omitempty"`
								AuthorUsername    string    `json:"author_username,omitempty"`
								AuthorName        string    `json:"author_name,omitempty"`
								CreatedAt         time.Time `json:"created_at,omitempty"`
								Lang              string    `json:"lang,omitempty"`
								Source            string    `json:"source,omitempty"`
								PossiblySensitive bool      `json:"possibly_sensitive,omitempty"`
							}
							out := prettyTweet{
								ID:                tweet.ID,
								Text:              tweet.Text,
								AuthorID:          tweet.AuthorID,
								CreatedAt:         tweet.CreatedAt,
								Lang:              tweet.Lang,
								Source:            tweet.Source,
								PossiblySensitive: tweet.PossiblySensitive,
							}
							if showUser && len(includes.Users) > 0 {
								out.AuthorUsername = includes.Users[0].Username
								out.AuthorName = includes.Users[0].Name
							}
							payload, err := json.MarshalIndent(out, "", "  ")
							if err != nil {
								return err
							}
							fmt.Println(string(payload))
							return nil
						}

						type rawEvent struct {
							Data     stream.StreamTweet    `json:"data"`
							Includes stream.StreamIncludes `json:"includes,omitempty"`
						}
						payload, err := json.Marshal(rawEvent{Data: tweet, Includes: includes})
						if err != nil {
							return err
						}
						fmt.Println(string(payload))
						return nil
					}

					// Display tweet (human output)
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

				if err == nil {
					lastDisconnect = "EOF"
				} else if errors.Is(err, io.EOF) {
					lastDisconnect = "EOF"
				} else if errors.Is(err, context.Canceled) {
					break
				} else {
					lastDisconnect = err.Error()
				}

				if err != nil && errors.Is(err, context.Canceled) {
					break
				}

				if ctx.Err() != nil {
					break
				}

				fmt.Fprintf(os.Stderr, "disconnected: %s\n", lastDisconnect)
				reconnects++

				wait := backoff
				if wait > maxBackoff {
					wait = maxBackoff
				}
				fmt.Fprintf(os.Stderr, "reconnecting in %s...\n", wait.Round(time.Second))

				timer := time.NewTimer(wait)
				select {
				case <-ctx.Done():
					timer.Stop()
					break
				case <-timer.C:
				}
				if ctx.Err() != nil {
					break
				}

				if backoff < maxBackoff {
					backoff *= 2
					if backoff > maxBackoff {
						backoff = maxBackoff
					}
				}
			}

			// Show summary
			duration := time.Since(startTime)
			fmt.Fprintf(os.Stderr, "\n\nStream summary: %s, %d tweets, reconnects=%d, last_disconnect=%s, last_ruleset=%s\n",
				duration.Round(time.Second), tweetCount, reconnects, lastDisconnect, lastRuleSet)
			if duration.Seconds() > 0 {
				rate := float64(tweetCount) / duration.Seconds() * 60
				fmt.Fprintf(os.Stderr, "Rate: %.1f tweets/minute\n", rate)
			}

			return nil
		},
	}

	cmd.Flags().StringArrayVar(&keywords, "keyword", nil, "Keyword to watch for (can be specified multiple times)")
	cmd.Flags().BoolVar(&autoSetup, "auto-setup", false, "Automatically set up stream rules for keywords")
	cmd.Flags().BoolVar(&showUser, "show-user", false, "Show author information")
	cmd.Flags().BoolVar(&showMeta, "show-meta", false, "Show additional metadata")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output newline-delimited JSON events")

	return cmd
}
