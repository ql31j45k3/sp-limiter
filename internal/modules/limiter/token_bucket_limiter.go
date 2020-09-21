package limiter

import (
	"github.com/ql31j45k3/sp-limiter/internal/utils/tools"
	"sync"
	"time"
)

func newTokenBucketLimiter(interval time.Duration, capacity int64) *tokenBucketLimiter {
	l := &tokenBucketLimiter{
		interval: interval,
		capacity: capacity,

		ip2Token: make(map[string]int64, capacity),
	}

	go func(l *tokenBucketLimiter) {
		ticker := time.NewTicker(l.interval)

		for {
			<-ticker.C
			l.addToken()
		}
	}(l)

	return l
}

// tokenBucketLimiter
type tokenBucketLimiter struct {
	interval time.Duration
	mu       sync.Mutex

	capacity int64
	ip2Token map[string]int64
}

func (l *tokenBucketLimiter) addToken() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for ip, _ := range l.ip2Token {
		l.ip2Token[ip] = l.capacity
	}
}

func (l *tokenBucketLimiter) TakeAvailable(ip string, block bool) (bool, int64) {
	if tools.IsEmpty(ip) {
		return false, 0
	}

	l.isExist(ip)

	blockFunc := func(l *tokenBucketLimiter) (bool, int64) {
		l.mu.Lock()

		tokenCount := l.ip2Token[ip]
		isTakeToken := (tokenCount - 1) >= 0
		if isTakeToken {
			defer l.mu.Unlock()
			tokenCount = tokenCount - 1
			l.ip2Token[ip] = tokenCount

			return true, l.capacity - tokenCount
		}


		l.mu.Unlock()
		time.Sleep(l.interval)

		l.mu.Lock()
		defer l.mu.Unlock()

		tokenCount = l.ip2Token[ip]
		if (tokenCount - 1) < 0 {
			return false, 0
		}

		tokenCount = tokenCount - 1
		l.ip2Token[ip] = tokenCount

		return true, l.capacity - tokenCount
	}

	nonBlockFunc := func(l *tokenBucketLimiter) (bool, int64) {
		l.mu.Lock()
		defer l.mu.Unlock()

		tokenCount := l.ip2Token[ip]
		if (tokenCount - 1) < 0 {
			return false, 0
		}

		tokenCount = tokenCount - 1
		l.ip2Token[ip] = tokenCount

		return true, l.capacity - tokenCount
	}

	if block {
		return blockFunc(l)
	}
	return nonBlockFunc(l)
}

func (l *tokenBucketLimiter) isExist(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.ip2Token[ip]; ok {
		return
	}

	// 第一次 IP 請求，初始化 token
	l.ip2Token[ip] = l.capacity
}
