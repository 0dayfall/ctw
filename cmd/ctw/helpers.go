package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/0dayfall/ctw/internal/client"
)

func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(os.Stdout, string(data))
	return err
}

func printRateLimits(snapshot client.RateLimitSnapshot) {
	if snapshot.Limit < 0 && snapshot.Remaining < 0 && snapshot.Reset < 0 {
		return
	}
	fmt.Fprintf(os.Stderr, "rate-limit limit=%d remaining=%d reset=%d\n", snapshot.Limit, snapshot.Remaining, snapshot.Reset)
}

func parseKeyValuePairs(pairs []string) (map[string]string, error) {
	values := make(map[string]string)
	for _, pair := range pairs {
		if !strings.Contains(pair, "=") {
			return nil, errors.New("expected key=value format")
		}
		parts := strings.SplitN(pair, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, errors.New("empty key in key=value pair")
		}
		values[key] = value
	}
	return values, nil
}
