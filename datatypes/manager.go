package datatypes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
)

type RateManager struct {
	limiterMemory sync.Map
	config        *RateLimitConfig
}

// LoadFromConfig creates a new manager from configuration
func LoadFromConfig(pathFile string) *RateManager {
	manager := NewRateManager()
	config, err := LoadConfig(pathFile)
	if err != nil {
		panic("missing config file")
	}
	manager.UseConfig(config)
	return manager
}

// NewRateManager creates a new in memory dict to store session data
// NOTE: unlimited map storage can lead to DoS attacks
func NewRateManager() *RateManager {
	return &RateManager{limiterMemory: sync.Map{}}
}

func (mngr *RateManager) getItem(id string) (*TokenBucketRateLimit, bool) {
	item, ok := mngr.limiterMemory.Load(id)
	if ok {
		limiter, ok2 := item.(*TokenBucketRateLimit)
		if ok2 {
			return limiter, true
		}
	}
	return nil, false
}

func (mngr *RateManager) putItem(id string, data *TokenBucketRateLimit) {
	mngr.limiterMemory.Store(id, data)
}

func (mngr *RateManager) UseConfig(config *RateLimitConfig) {
	mngr.config = config
}

// BuildFor search for a declared rate limit configuration and returns the middleware with that configuration set
func (mngr *RateManager) BuildFor(methodType string, pathPattern string) gin.HandlerFunc {
	endpointConfig := mngr.findConfig(methodType, pathPattern)
	if endpointConfig == nil {
		// we do not create a middleware
		return func(context *gin.Context) {
			context.Next()
		}
	}
	// we return a rate limit middleware for that endpoint
	return mngr.createMiddleware(endpointConfig)
}

// CheckStatus returns the status information of given endpoint for a caller
func (mngr *RateManager) CheckStatus(ctx *gin.Context, methodType string, pathPattern string) (*TokenBucketRateLimit, bool) {
	limiterId := fmt.Sprintf("%s-%s %s", ctx.ClientIP(), methodType, pathPattern)
	status, found := mngr.getItem(limiterId)
	return status, found
}

// findConfig search for matching configuration
func (mngr *RateManager) findConfig(methodType string, pathPattern string) *RateLimitsPerEndpoint {
	for _, item := range mngr.config.RateLimitsPerEndpoint {
		chunks := strings.Split(item.Endpoint, " ")
		if len(chunks) == 2 {
			method := chunks[0]
			pattern := chunks[1]
			if method == methodType && pattern == pathPattern {
				// return
				return &item
			}
		}
	}
	return nil
}

func (mngr *RateManager) createMiddleware(endpointInfo *RateLimitsPerEndpoint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. find the caller client
		limiterId := fmt.Sprintf("%s-%s", ctx.ClientIP(), endpointInfo.Endpoint)
		status, found := mngr.getItem(limiterId)
		if !found {
			// this is first time client is calling this endpoint
			// build new instance and use it
			endpointLimiter := NewTokenBucketRateLimit(int32(endpointInfo.Burst), int32(endpointInfo.Sustained))
			mngr.putItem(limiterId, endpointLimiter)
			status = endpointLimiter
		}
		_, statusCode, statusErr := status.IsAllowed()
		if statusErr != nil {
			ctx.AbortWithError(statusCode, statusErr)
			return
		}
		ctx.Next()
	}
}
