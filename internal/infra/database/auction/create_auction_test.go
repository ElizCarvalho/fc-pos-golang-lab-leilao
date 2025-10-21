package auction

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAuctionDuration(t *testing.T) {
	tests := []struct {
		name           string
		envValue       string
		expectedResult time.Duration
	}{
		{
			name:           "Valid duration from env",
			envValue:       "10s",
			expectedResult: 10 * time.Second,
		},
		{
			name:           "Invalid duration, should use default",
			envValue:       "invalid",
			expectedResult: 5 * time.Minute,
		},
		{
			name:           "Empty env, should use default",
			envValue:       "",
			expectedResult: 5 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			os.Setenv("AUCTION_DURATION", tt.envValue)
			defer os.Unsetenv("AUCTION_DURATION")

			result := getAuctionDuration()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGetAuctionCheckInterval(t *testing.T) {
	tests := []struct {
		name           string
		envValue       string
		expectedResult time.Duration
	}{
		{
			name:           "Valid interval from env",
			envValue:       "30s",
			expectedResult: 30 * time.Second,
		},
		{
			name:           "Invalid interval, should use default",
			envValue:       "invalid",
			expectedResult: 1 * time.Minute,
		},
		{
			name:           "Empty env, should use default",
			envValue:       "",
			expectedResult: 1 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			os.Setenv("AUCTION_CHECK_INTERVAL", tt.envValue)
			defer os.Unsetenv("AUCTION_CHECK_INTERVAL")

			result := getAuctionCheckInterval()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestUpdateAuctionStatus(t *testing.T) {
	t.Run("valid duration", func(t *testing.T) {
		os.Setenv("AUCTION_DURATION", "10s")
		defer os.Unsetenv("AUCTION_DURATION")

		result := getAuctionDuration()
		assert.Equal(t, 10*time.Second, result)
	})

	t.Run("invalid duration", func(t *testing.T) {
		os.Setenv("AUCTION_DURATION", "invalid")
		defer os.Unsetenv("AUCTION_DURATION")

		result := getAuctionDuration()
		assert.Equal(t, 5*time.Minute, result)
	})
}

func TestCreateAuction(t *testing.T) {
	t.Run("valid check interval", func(t *testing.T) {
		os.Setenv("AUCTION_CHECK_INTERVAL", "30s")
		defer os.Unsetenv("AUCTION_CHECK_INTERVAL")

		result := getAuctionCheckInterval()
		assert.Equal(t, 30*time.Second, result)
	})

	t.Run("invalid check interval", func(t *testing.T) {
		os.Setenv("AUCTION_CHECK_INTERVAL", "invalid")
		defer os.Unsetenv("AUCTION_CHECK_INTERVAL")

		result := getAuctionCheckInterval()
		assert.Equal(t, 1*time.Minute, result)
	})
}

func TestStartAutoCloseRoutine(t *testing.T) {
	t.Run("auto close configuration", func(t *testing.T) {
		os.Setenv("AUCTION_DURATION", "100ms")
		os.Setenv("AUCTION_CHECK_INTERVAL", "50ms")
		defer func() {
			os.Unsetenv("AUCTION_DURATION")
			os.Unsetenv("AUCTION_CHECK_INTERVAL")
		}()

		duration := getAuctionDuration()
		interval := getAuctionCheckInterval()

		assert.Equal(t, 100*time.Millisecond, duration)
		assert.Equal(t, 50*time.Millisecond, interval)
	})
}
