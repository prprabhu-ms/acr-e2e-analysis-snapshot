package api

import (
	"log"
	"strings"
	"time"

	"github.com/cli/go-gh"
)

type retryManager struct {
	count5XX                        int
	countSecondaryRateLimitExceeded int
}

func (m *retryManager) retryAfter(err error) (bool, time.Duration) {
	if isError5XX(err) {
		m.count5XX++
		backoff := time.Duration(m.count5XX) * baseBackoff5XX
		if backoff > maxBackoff {
			return false, 0
		}
		return true, backoff
	}

	// A user has exceeded their rate limit, more info: https://docs.github.com/en/rest/overview/resources-in-the-rest-api#rate-limiting
	if isUserRateLimitExceeded(err) {
		backoff := time.Duration(getUserRateLimitResetTime()-time.Now().Unix()) * time.Second
		return true, backoff
	}

	// A specific API is called too often, more info: https://docs.github.com/en/rest/overview/resources-in-the-rest-api#secondary-rate-limits
	if isSecondaryRateLimitExceeded(err) {
		m.countSecondaryRateLimitExceeded++
		// TODO: instead of an exponential backoff, we should use the `retry-after` response header
		backoff := time.Duration(m.countSecondaryRateLimitExceeded) * baseBackoffSecondaryRateLimit
		if backoff > maxBackoff {
			return false, 0
		}
		return true, backoff
	}

	return false, 0
}

const baseBackoff5XX = 1 * time.Minute
const baseBackoffSecondaryRateLimit = 3 * time.Minute
const maxBackoff = 25 * time.Minute

func isError5XX(err error) bool {
	return strings.Contains(err.Error(), "HTTP 500") || strings.Contains(err.Error(), "HTTP 502") || strings.Contains(err.Error(), "HTTP 503")
}

func isUserRateLimitExceeded(err error) bool {
	return strings.Contains(err.Error(), "API rate limit exceeded for user ID")
}

func getUserRateLimitResetTime() int64 {
	client, err := gh.RESTClient(nil)

	if err != nil {
		log.Fatal(err)
	}

	response := &struct {
		Resources struct {
			Core struct {
				Reset int64
			}
		}
	}{}

	client.Get("rate_limit", response)
	return response.Resources.Core.Reset
}

func isSecondaryRateLimitExceeded(err error) bool {
	return strings.Contains(err.Error(), "You have exceeded a secondary rate limit")
}
