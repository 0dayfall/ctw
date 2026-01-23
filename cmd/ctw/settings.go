package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/0dayfall/ctw/internal/config"
	"github.com/spf13/cobra"
)

type Settings struct {
	BaseURL          string
	BearerToken      string
	UserAgent        string
	Timeout          time.Duration
	Retry            int
	PrettyOutput     bool
	StreamBackoffMax time.Duration
	ConfigPath       string
	ConfigLoaded     bool
}

var (
	configPathFlag string
	timeoutFlag    time.Duration
	retryFlag      int
	prettyFlag     bool

	resolvedSettings Settings
	settingsLoaded   bool
	prettyOutput     bool
)

func ensureSettings(cmd *cobra.Command) error {
	if settingsLoaded {
		return nil
	}

	cfgPath := strings.TrimSpace(configPathFlag)
	if cfgPath == "" {
		path, err := config.DefaultPath()
		if err != nil {
			return err
		}
		cfgPath = path
	}

	cfg, loaded, err := config.Load(cfgPath)
	if err != nil {
		return err
	}
	if cmd != nil && cmd.Flags().Changed("config") && !loaded {
		return fmt.Errorf("config file not found: %s", cfgPath)
	}

	settings := Settings{
		BaseURL:          strings.TrimSpace(baseURLFlag),
		BearerToken:      strings.TrimSpace(cfg.Auth.BearerToken),
		UserAgent:        strings.TrimSpace(cfg.HTTP.UserAgent),
		Timeout:          cfg.HTTP.Timeout.Std(),
		Retry:            cfg.HTTP.Retry,
		PrettyOutput:     cfg.Output.Pretty,
		StreamBackoffMax: cfg.Stream.BackoffMax.Std(),
		ConfigPath:       cfgPath,
		ConfigLoaded:     loaded,
	}

	if err := applyEnvOverrides(&settings); err != nil {
		return err
	}
	if err := applyFlagOverrides(cmd, &settings); err != nil {
		return err
	}

	resolvedSettings = settings
	settingsLoaded = true
	prettyOutput = settings.PrettyOutput
	return nil
}

func applyEnvOverrides(settings *Settings) error {
	if value := strings.TrimSpace(os.Getenv("BEARER_TOKEN")); value != "" {
		settings.BearerToken = value
	}
	if value := strings.TrimSpace(os.Getenv("USER_AGENT")); value != "" {
		settings.UserAgent = value
	}
	if value := strings.TrimSpace(os.Getenv("CTW_BASE_URL")); value != "" {
		settings.BaseURL = value
	}
	if value := strings.TrimSpace(os.Getenv("CTW_TIMEOUT")); value != "" {
		timeout, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("parse CTW_TIMEOUT: %w", err)
		}
		settings.Timeout = timeout
	}
	if value := strings.TrimSpace(os.Getenv("CTW_RETRY")); value != "" {
		retry, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("parse CTW_RETRY: %w", err)
		}
		settings.Retry = retry
	}
	if value := strings.TrimSpace(os.Getenv("CTW_PRETTY")); value != "" {
		pretty, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("parse CTW_PRETTY: %w", err)
		}
		settings.PrettyOutput = pretty
	}
	if value := strings.TrimSpace(os.Getenv("CTW_STREAM_BACKOFF_MAX")); value != "" {
		backoff, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("parse CTW_STREAM_BACKOFF_MAX: %w", err)
		}
		settings.StreamBackoffMax = backoff
	}
	return nil
}

func applyFlagOverrides(cmd *cobra.Command, settings *Settings) error {
	if cmd == nil {
		return nil
	}

	if cmd.Flags().Changed("base-url") {
		settings.BaseURL = strings.TrimSpace(baseURLFlag)
	}
	if cmd.Flags().Changed("bearer-token") {
		settings.BearerToken = strings.TrimSpace(bearerTokenFlag)
	}
	if cmd.Flags().Changed("user-agent") {
		settings.UserAgent = strings.TrimSpace(userAgentFlag)
	}
	if cmd.Flags().Changed("timeout") {
		settings.Timeout = timeoutFlag
	}
	if cmd.Flags().Changed("retry") {
		settings.Retry = retryFlag
	}
	if cmd.Flags().Changed("pretty") {
		settings.PrettyOutput = prettyFlag
	}

	return nil
}
