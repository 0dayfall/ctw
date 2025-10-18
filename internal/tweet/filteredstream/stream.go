package tweet

import (
	"context"

	"github.com/0dayfall/ctw/internal/client"
)

// Stream is a convenience wrapper that constructs a Service on the fly and
// proxies to its Stream method. New code should prefer using a shared Service
// instance directly so that HTTP client configuration can be reused.
func Stream(ctx context.Context, c *client.Client, fields map[string]string) (StreamEnvelope, client.RateLimitSnapshot, error) {
	return NewService(c).Stream(ctx, fields)
}
