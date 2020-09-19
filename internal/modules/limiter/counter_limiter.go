package limiter

import (
	"github.com/ql31j45k3/sp-limiter/internal/utils/tools"
	"strconv"
	"sync"
	"time"
)

func newCounterLimit(interval time.Duration, maxCount int) *counterLimit {
	l := &counterLimit{
		interval: interval,
		maxCount: maxCount,

		ip2count: make(map[string]int),
	}

	go func(l *counterLimit) {
		ticker := time.NewTicker(l.interval)

		for {
			<-ticker.C
			l.Zero()
		}
	}(l)

	return l
}

type counterLimit struct {
	interval time.Duration
	mu       sync.Mutex

	maxCount int
	ip2count map[string]int
}

func (l *counterLimit) Increase(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if tools.IsEmpty(ip) {
		return
	}

	l.ip2count[ip] += 1
}

func (l *counterLimit) Zero() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.ip2count = make(map[string]int)
}

func (l *counterLimit) IsAvailable(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if tools.IsEmpty(ip) {
		return false
	}

	return l.ip2count[ip] < l.maxCount
}

func (l *counterLimit) GetCount(ip string) string {
	l.mu.Lock()
	defer l.mu.Unlock()

	if tools.IsEmpty(ip) {
		return ""
	}

	return strconv.Itoa(l.ip2count[ip])
}
