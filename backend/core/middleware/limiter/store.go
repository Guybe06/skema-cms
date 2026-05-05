package limiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type ipStore struct {
	mu       sync.RWMutex
	limiters map[string]*ipLimiter
	rps      rate.Limit
	burst    int
}

var (
	globalStore = newIPStore(RateLimitGlobal, RateLimitBurst*10)
	authStore   = newIPStore(RateLimitAuth, RateLimitBurst)
)

func newIPStore(rpm, burst int) *ipStore {
	s := &ipStore{
		limiters: make(map[string]*ipLimiter),
		rps:      rate.Limit(rpm) / 60,
		burst:    burst,
	}
	go s.cleanup()
	return s
}

func (s *ipStore) get(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry, ok := s.limiters[ip]; ok {
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	l := rate.NewLimiter(s.rps, s.burst)
	s.limiters[ip] = &ipLimiter{limiter: l, lastSeen: time.Now()}
	return l
}

func (s *ipStore) cleanup() {
	ticker := time.NewTicker(CleanupIntervalMin * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		cutoff := time.Now().Add(-LimiterTTLMin * time.Minute)
		for ip, entry := range s.limiters {
			if entry.lastSeen.Before(cutoff) {
				delete(s.limiters, ip)
			}
		}
		s.mu.Unlock()
	}
}
