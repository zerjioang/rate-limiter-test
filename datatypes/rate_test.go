package datatypes

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestNewTokenBucketRateLimit(t *testing.T) {
	t.Run("create-bucket", func(t *testing.T) {
		limiter := NewTokenBucketRateLimit(100, 50)
		assert.NotNil(t, limiter)

		assert.Equal(t, limiter.maxRate, int32(100))
		assert.Equal(t, limiter.RequestPerMinute, int32(50))
	})

	t.Run("is-allowed-once", func(t *testing.T) {
		limiter := NewTokenBucketRateLimit(100, 50)
		assert.NotNil(t, limiter)

		avail, status, err := limiter.IsAllowed()
		fmt.Println(avail, status, err)
		assert.Equal(t, avail, int32(49))
		assert.Equal(t, status, http.StatusOK)
		assert.NoError(t, err)
	})

	t.Run("all-tokens-used", func(t *testing.T) {
		limiter := NewTokenBucketRateLimit(100, 50)
		assert.NotNil(t, limiter)

		availableExpected := 50
		for i := 0; i < 50; i++ {
			avail, status, err := limiter.IsAllowed()
			availableExpected--
			fmt.Println(avail, status, err)
			assert.Equal(t, avail, int32(availableExpected))
			assert.Equal(t, status, http.StatusOK)
			assert.NoError(t, err)
		}
	})

	t.Run("all-tokens-used", func(t *testing.T) {
		limiter := NewTokenBucketRateLimit(100, 50)
		assert.NotNil(t, limiter)

		avail, status, err := limiter.IsAllowed()
		availableExpected := 50
		for i := 0; i < 50; i++ {
			time.Sleep(1 * time.Second)
			avail, status, err = limiter.IsAllowed()
			availableExpected--
			fmt.Println(avail, status, err)
			assert.Equal(t, avail, int32(availableExpected))
			assert.Equal(t, status, http.StatusOK)
			assert.NoError(t, err)
		}
	})

	t.Run("run-out-of-tokens", func(t *testing.T) {
		limiter := NewTokenBucketRateLimit(100, 50)
		assert.NotNil(t, limiter)

		availableExpected := 50
		for i := 0; i < 50; i++ {
			avail, status, err := limiter.IsAllowed()
			availableExpected--
			assert.Equal(t, avail, int32(availableExpected))
			assert.Equal(t, status, http.StatusOK)
			assert.NoError(t, err)
		}

		// next request should be blocked
		avail, status, err := limiter.IsAllowed()
		assert.Equal(t, avail, int32(0))
		assert.Equal(t, status, http.StatusBadRequest)
		assert.Error(t, err)
	})
}
