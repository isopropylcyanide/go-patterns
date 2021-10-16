package rate_limiter

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// Zen: Any sane production system would try to implement a rate limits on its resource
// to prevent a cascading failure or misuse of the resources. Mostly, it is implemented
// using a leaky token bucket algorithm that exposes two parameters, a burst rate and a
// replenishment rate. A burst rate measures how many requests can be made when the bucket

// Generally, implemented at the server (because clients may bypass), it can be done at
// the client layer as an optimization

// Here we simulate a dummy client operation that pretends to be costly and needs to be
// rate limited. The different combinations of the rate limiter are provided through tests.
type apiConnection struct {
	rateLimiter *rate.Limiter
}

func Open(replenishAfter time.Duration, burst int) *apiConnection {
	return &apiConnection{
		rateLimiter: rate.NewLimiter(rate.Every(replenishAfter), burst),
	}
}

// ReadFile Random method that is rate limited. The request goes over the wire,
// and we may want to cancel. Hence, context is passed in the request.
func (a *apiConnection) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	return nil
}
