package api

import (
	"strings"
	"time"
)

type retryManager struct {
	count5XX      int
	countThrottle int
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

	if isErrorThrottle(err) {
		m.countThrottle++
		backoff := time.Duration(m.countThrottle) * baseBackoffThrottle
		if backoff > maxBackoff {
			return false, 0
		}
		return true, backoff
	}
	return false, 0
}

const baseBackoff5XX = 1 * time.Minute
const baseBackoffThrottle = 10 * time.Minute
const maxBackoff = 25 * time.Minute

func isError5XX(err error) bool {
	return strings.Contains(err.Error(), "HTTP 500") || strings.Contains(err.Error(), "HTTP 502") || strings.Contains(err.Error(), "HTTP 503")
}

func isErrorThrottle(err error) bool {
	return strings.Contains(err.Error(), "API rate limit exceeded")
}
