package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/0dayfall/ctw/internal/config"
	lookupsvc "github.com/0dayfall/ctw/internal/users/lookup"
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

type rawConfig struct {
	Auth struct {
		BearerToken string `toml:"bearer_token"`
	} `toml:"auth"`
}

func init() {
	rootCmd.AddCommand(newInitCommand())
}

func newInitCommand() *cobra.Command {
	var (
		writeConfig    bool
		noWrite        bool
		forceOverwrite bool
		allowPlaintext bool
		username       string
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize ctw configuration and verify API access",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureSettings(cmd); err != nil {
				return err
			}

			if noWrite {
				writeConfig = false
			}

			if writeConfig && noWrite {
				return errors.New("cannot combine --write and --no-write")
			}

			token, source, err := resolveBearerTokenSource()
			if err != nil {
				return err
			}
			if token == "" {
				return errors.New("bearer token not found (use --bearer-token, BEARER_TOKEN, or config file)")
			}

			if username == "" {
				username = "twitter"
			}

			client, err := newClientFromFlags()
			if err != nil {
				return err
			}

			service := lookupsvc.NewService(client)
			user, rateLimits, err := service.LookupUsername(cmd.Context(), username, nil)
			if err != nil {
				printRateLimits(rateLimits)
				return err
			}

			fmt.Printf("✔ Config path: %s\n", resolvedSettings.ConfigPath)
			if source == "" {
				source = "unknown"
			}
			fmt.Printf("✔ Bearer token: found (from %s)\n", source)
			fmt.Printf("✔ Auth: OK (user: @%s, id: %s)\n", user.UserName, user.ID)
			fmt.Printf("✔ API access: OK\n")

			if writeConfig {
				if err := writeConfigFile(token, source, allowPlaintext, forceOverwrite); err != nil {
					return err
				}
			} else {
				fmt.Printf("✔ Wrote config: skipped\n")
			}

			fmt.Printf("Done. Try: ctw search recent --query \"golang\"\n")

			if source == "env:BEARER_TOKEN" || (source == "flag" && os.Getenv("BEARER_TOKEN") == "") {
				fmt.Printf("Add this to your shell rc:\n")
				fmt.Printf("export BEARER_TOKEN=\"%s\"\n", token)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&writeConfig, "write", true, "Write config file (default)")
	cmd.Flags().BoolVar(&noWrite, "no-write", false, "Do not write a config file")
	cmd.Flags().BoolVar(&forceOverwrite, "force", false, "Overwrite existing config file")
	cmd.Flags().BoolVar(&allowPlaintext, "allow-plaintext", false, "Allow writing raw bearer token to config")
	cmd.Flags().StringVar(&username, "username", "twitter", "Username to verify API access")

	return cmd
}

func resolveBearerTokenSource() (string, string, error) {
	if bearerTokenFlag != "" {
		return strings.TrimSpace(bearerTokenFlag), "flag", nil
	}

	if envToken := strings.TrimSpace(os.Getenv("BEARER_TOKEN")); envToken != "" {
		return envToken, "env:BEARER_TOKEN", nil
	}

	if resolvedSettings.ConfigLoaded {
		rawToken := loadRawBearerToken(resolvedSettings.ConfigPath)
		if strings.HasPrefix(rawToken, "env:") {
			envName := strings.TrimPrefix(rawToken, "env:")
			envValue := strings.TrimSpace(os.Getenv(envName))
			if envValue != "" {
				return envValue, "env:" + envName, nil
			}
		}
		if strings.TrimSpace(rawToken) != "" {
			return strings.TrimSpace(rawToken), "config", nil
		}
	}

	return "", "", nil
}

func loadRawBearerToken(path string) string {
	if path == "" {
		return ""
	}
	var raw rawConfig
	if _, err := toml.DecodeFile(path, &raw); err != nil {
		return ""
	}
	return strings.TrimSpace(raw.Auth.BearerToken)
}

func writeConfigFile(token, source string, allowPlaintext, force bool) error {
	path := resolvedSettings.ConfigPath
	if path == "" {
		var err error
		path, err = config.DefaultPath()
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(path); err == nil && !force {
		fmt.Printf("✔ Wrote config: skipped (exists; use --force to overwrite)\n")
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	tokenValue := ""
	if strings.HasPrefix(source, "env:") {
		tokenValue = source
	} else if allowPlaintext {
		tokenValue = token
	} else {
		tokenValue = "env:BEARER_TOKEN"
	}

	userAgent := strings.TrimSpace(resolvedSettings.UserAgent)
	if userAgent == "" {
		userAgent = "ctw"
	}

	content := fmt.Sprintf(`[auth]
bearer_token = "%s"

[http]
user_agent = "%s"
timeout = "%s"
retry = %d

[output]
pretty = %t

[stream]
backoff_max = "%s"
`, tokenValue, userAgent, resolvedSettings.Timeout, resolvedSettings.Retry, resolvedSettings.PrettyOutput, resolvedSettings.StreamBackoffMax)

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return err
	}

	fmt.Printf("✔ Wrote config: %s\n", path)
	return nil
}
