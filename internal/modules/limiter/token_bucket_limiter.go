package limiter

import (
	"sync"
	"time"

	"github.com/ql31j45k3/sp-limiter/internal/utils/tools"
)

// newTokenBucketLimiter 初始化參數
// 同時執行一個 goroutine 背景執行 ticker 依照 interval 參數，觸發重新置入 token
func newTokenBucketLimiter(interval time.Duration, capacity int64) *tokenBucketLimiter {
	l := &tokenBucketLimiter{
		interval: interval,
		capacity: capacity,

		ip2token: make(map[string]int64),
	}

	go func(l *tokenBucketLimiter) {
		ticker := time.NewTicker(l.interval)

		for {
			<-ticker.C
			// 令牌演算法邏輯上要先產生 token，但此功能是每個 IP 設置 token，
			// 故初始 token 邏輯在第一次取得 IP 時建立 token，在每次觸發 ticker 做重置 map 邏輯
			l.Zero()
		}
	}(l)

	return l
}

type tokenBucketLimiter struct {
	interval time.Duration
	mu       sync.Mutex

	capacity int64
	ip2token map[string]int64
}

// Zero 將紀錄 ip 對應的 token 做初始化
func (l *tokenBucketLimiter) Zero() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.ip2token = make(map[string]int64)
}

// TakeAvailable 確認是否可用 (尚未到達限流條件)，並同時減少一個 token
// block 控制是否使用阻塞方式，等待 token (只會阻塞一次，第二次在未取得 token 回傳失敗)
func (l *tokenBucketLimiter) TakeAvailable(ip string, block bool) (bool, int64) {
	if tools.IsEmpty(ip) {
		return false, 0
	}

	l.isExist(ip)

	// 處理如果 token 已沒有，再次等待取 token 邏輯（只處理一次重新等待，並非一定要等待拿到 token 為止）
	blockFunc := func(l *tokenBucketLimiter) (bool, int64) {
		l.mu.Lock()

		// 第一次確認是否可以取得 token
		tokenCount := l.ip2token[ip]
		// 驗證是否可取得 token，最多減到為 0 的數字
		isTakeToken := (tokenCount - 1) >= 0
		if isTakeToken {
			// 成功取得 token 加上 defer Unlock
			defer l.mu.Unlock()

			tokenCount = tokenCount - 1
			l.ip2token[ip] = tokenCount

			return true, l.capacity - tokenCount
		}

		// 未成功取得 token，馬上 Unlock，並免佔用鎖 無法執行 addToken 邏輯
		l.mu.Unlock()
		// 用 sleep 方式做等待 token 產生邏輯
		time.Sleep(l.interval)

		// 第二次確認，是否可取得 token，並需重新上鎖一次
		l.mu.Lock()
		defer l.mu.Unlock()

		tokenCount = l.ip2token[ip]
		// 計算為負數，代表 token 取得失敗
		if (tokenCount - 1) < 0 {
			return false, 0
		}

		tokenCount = tokenCount - 1
		l.ip2token[ip] = tokenCount

		return true, l.capacity - tokenCount
	}

	// 處理如果沒 token 馬上返回，不重新等待取 token 流程
	nonBlockFunc := func(l *tokenBucketLimiter) (bool, int64) {
		l.mu.Lock()
		defer l.mu.Unlock()

		tokenCount := l.ip2token[ip]
		// 計算為負數，代表 token 取得失敗
		if (tokenCount - 1) < 0 {
			return false, 0
		}

		tokenCount = tokenCount - 1
		l.ip2token[ip] = tokenCount

		return true, l.capacity - tokenCount
	}

	if block {
		return blockFunc(l)
	}
	return nonBlockFunc(l)
}

// isExist 確認是否已經存在此 IP，未存在做第一次初始化 token 邏輯
func (l *tokenBucketLimiter) isExist(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.ip2token[ip]; ok {
		return
	}

	// 第一次 IP 請求，初始化 token
	l.ip2token[ip] = l.capacity
}
