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

		ip2Token:          make(map[string]chan struct{}, capacity),
		ip2AvailableToken: make(map[string]int64),
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

// tokenBucketLimiter，用 channel 模擬 token bucket
type tokenBucketLimiter struct {
	interval time.Duration
	mu       sync.Mutex

	capacity          int64
	ip2Token          map[string]chan struct{}
	ip2AvailableToken map[string]int64
}

func (l *tokenBucketLimiter) addToken() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for k, v := range l.ip2Token {
		ip := k
		var i int64
		for i = 0; i < l.capacity; i++ {
			select {
			case v <- struct{}{}:
			default:
				// 代表容量已滿，後續直接略過
				break
			}
		}

		l.ip2AvailableToken[ip] = 0
	}
}

func (l *tokenBucketLimiter) TakeAvailable(ip string, block bool) (bool, int64) {
	if tools.IsEmpty(ip) {
		return false, 0
	}

	// 目前數據結構，IP 對應 chan，由於 IP 無法事先知道，先用惰性方式，第一次再置入 token
	l.isExist(ip)

	l.mu.Lock()
	if !block {
		// 客戶端沒有阻塞等待 token 才可使用 defer 方式解除 lock
		defer l.mu.Unlock()
	}

	// 選擇等待拿到 token
	if block {
		// 要解除 lock，才可等待 token 新增
		l.mu.Unlock()
		select {
		case <-l.ip2Token[ip]:
			l.ip2AvailableToken[ip] += 1
			return true, l.ip2AvailableToken[ip]
		}
	}

	// 選擇沒有 token，直接 return
	select {
	case <-l.ip2Token[ip]:
		l.ip2AvailableToken[ip] += 1
		return true, l.ip2AvailableToken[ip]
	default:
		return false, 0
	}
}

func (l *tokenBucketLimiter) isExist(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.ip2Token[ip]; ok {
		return
	}

	// 第一次 IP 請求，初始化 token
	l.ip2Token[ip] = make(chan struct{}, l.capacity)
	v := l.ip2Token[ip]
	var i int64
	for i = 0; i < l.capacity; i++ {
		select {
		case v <- struct{}{}:
		}
	}
}
