package main

import (
	"fmt"
	"os"

	"github.com/0dayfall/ctw/internal/client"
	"github.com/spf13/cobra"
)

var (
	bearerTokenFlag string
	baseURLFlag     string
	userAgentFlag   string
)

var rootCmd = &cobra.Command{
	Use:   "ctw",
	Short: "Twitter v2 command line client",
	Long:  "ctw is a command line interface for interacting with selected Twitter v2 API endpoints.",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&bearerTokenFlag, "bearer-token", "", "Twitter API bearer token (defaults to BEARER_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&baseURLFlag, "base-url", "", "Override API base URL (defaults to https://api.twitter.com/)")
	rootCmd.PersistentFlags().StringVar(&userAgentFlag, "user-agent", "", "Override HTTP User-Agent header")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newClientFromFlags() (*client.Client, error) {
	cfg := client.Config{
		BaseURL:     baseURLFlag,
		BearerToken: bearerTokenFlag,
		UserAgent:   userAgentFlag,
	}
	return client.New(cfg)
}
