package tweet

import (
	"context"

	"github.com/0dayfall/ctw/internal/client"
)

// AddRule is a convenience wrapper that delegates to the Service abstraction.
func AddRule(ctx context.Context, c *client.Client, cmd AddCommand, dryRun bool) (RulesResponse, client.RateLimitSnapshot, error) {
	return NewService(c).AddRule(ctx, cmd, dryRun)
}

// GetRules is a convenience wrapper that delegates to the Service abstraction.
func GetRules(ctx context.Context, c *client.Client) (RulesResponse, client.RateLimitSnapshot, error) {
	return NewService(c).GetRules(ctx)
}
