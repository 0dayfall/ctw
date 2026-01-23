package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// Config captures user configurable settings for ctw.
type Config struct {
	Auth struct {
		BearerToken string `toml:"bearer_token"`
	} `toml:"auth"`

	HTTP struct {
		UserAgent string   `toml:"user_agent"`
		Timeout   Duration `toml:"timeout"`
		Retry     int      `toml:"retry"`
	} `toml:"http"`

	Output struct {
		Pretty bool `toml:"pretty"`
	} `toml:"output"`

	Stream struct {
		BackoffMax Duration `toml:"backoff_max"`
	} `toml:"stream"`
}

// Default returns a config populated with default values.
func Default() Config {
	var cfg Config
	cfg.HTTP.Timeout = Duration(15 * time.Second)
	cfg.HTTP.Retry = 3
	cfg.Output.Pretty = true
	cfg.Stream.BackoffMax = Duration(2 * time.Minute)
	return cfg
}

// DefaultPath returns the default config location for the current OS.
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		if appdata := os.Getenv("APPDATA"); appdata != "" {
			return filepath.Join(appdata, "ctw", "config.toml"), nil
		}
	}

	return filepath.Join(home, ".config", "ctw", "config.toml"), nil
}

// Load reads a config file from the provided path, applying defaults and env refs.
// The boolean return value indicates whether a file was loaded.
func Load(path string) (Config, bool, error) {
	cfg := Default()
	if path == "" {
		var err error
		path, err = DefaultPath()
		if err != nil {
			return cfg, false, err
		}
	}

	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, false, nil
		}
		return cfg, false, err
	}

	if _, err := toml.Decode(string(b), &cfg); err != nil {
		return cfg, false, err
	}

	resolveEnvRefs(&cfg)

	return cfg, true, nil
}

func resolveEnvRefs(cfg *Config) {
	cfg.Auth.BearerToken = expandEnvRef(cfg.Auth.BearerToken)
	cfg.HTTP.UserAgent = expandEnvRef(cfg.HTTP.UserAgent)
}

func expandEnvRef(value string) string {
	if strings.HasPrefix(value, "env:") {
		return os.Getenv(strings.TrimPrefix(value, "env:"))
	}
	return value
}

// Duration wraps time.Duration for TOML decoding (strings like "15s").
type Duration time.Duration

func (d *Duration) UnmarshalText(text []byte) error {
	value, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = Duration(value)
	return nil
}

func (d Duration) Std() time.Duration {
	return time.Duration(d)
}
