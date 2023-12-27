package datatypes

import (
	"errors"
	"net/http"
	"sync/atomic"
	"time"
)

type TokenBucketRateLimit struct {
	//
	// configuration related fields. they should not change at runtime
	//
	// max rate of tokens to have.
	// Example = 100
	maxRate int32
	// sustained requests per minute.
	// Also known as refresh rate of the tokens
	// Example = 50 req/min
	RequestPerMinute int32
	//
	// user state related fields
	//
	// current user (client) accumulated tokens
	// Max user tokens are 100
	currentTokens atomic.Int32
	// To automatically refresh user tokens, we also need to store client las call time
	LastRequestTime atomic.Int64
	// Conditions: on every http call, deduce 1 token
	// If no token left, deny the request
}

// NewTokenBucketRateLimit is a constructor like function that creates a new TokenBucketRateLimit
// with custom configuration for given endpoint
func NewTokenBucketRateLimit(burst int32, sustained int32) *TokenBucketRateLimit {
	limiter := &TokenBucketRateLimit{
		maxRate:          burst,
		currentTokens:    atomic.Int32{},
		RequestPerMinute: sustained,
		// initialize the rate limit with actual time
		LastRequestTime: atomic.Int64{},
	}
	limiter.currentTokens.Store(sustained)
	limiter.LastRequestTime.Store(time.Now().Unix())
	return limiter
}

// AvailableTokens returns the number of available tokens
func (limit *TokenBucketRateLimit) AvailableTokens() int32 {
	return limit.currentTokens.Load()
}

// fastRefill executes the refill logic as fast as possible on each request
func (limit *TokenBucketRateLimit) fastRefill() {
	// first, we verify if current user has refill available
	// to do that, we compute the delta time
	currentTime := time.Now().Unix()
	diffMs := (currentTime - limit.LastRequestTime.Load()) * 1000
	if diffMs > 0 {
		// compute the tokens to be added
		newTokensPerMs := float64(limit.RequestPerMinute) / 60.0 / 1000
		// precision is lost, but we are working with token units no..subunits like 0,01 ..etc
		newTokens := int32(float64(diffMs) * newTokensPerMs)
		available := limit.currentTokens.Load()
		// check if user has some tokens available to refill depending on the refill rate
		// at the end of the process, verify that max token values is still present
		if available+newTokens >= limit.maxRate {
			limit.currentTokens.Store(limit.maxRate)
		} else {
			// refill user tokens
			limit.currentTokens.Add(newTokens)
		}
	}
}

// IsAllowed returns nil if request is allowed to continue given rate limit conditions
// Returns error, if request not allowed and specifies the reason
// It also returns the remaining tokens for the user
func (limit *TokenBucketRateLimit) IsAllowed() (int32, int, error) {
	// first, we verify if current user has refill available, and how much
	limit.fastRefill()
	// then, we apply limit logic
	available := limit.currentTokens.Load()
	if available <= 0 {
		return 0, http.StatusBadRequest, errors.New("request rejected. Client has no tokens")
	}
	// If the limit for the specified route is not exceeded
	// the endpoint should consume one token
	limit.currentTokens.Add(-1)
	// store client last call time
	limit.LastRequestTime.Store(time.Now().Unix())
	return available - 1, http.StatusOK, nil
}
