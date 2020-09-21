package limiter

import (
	"strconv"
	"sync"
	"time"

	"github.com/ql31j45k3/sp-limiter/internal/utils/tools"
)

// newCounterLimit 初始化參數
// 同時執行一個 goroutine 背景執行 ticker 依照 interval 參數，觸發做重新解除限流邏輯
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

// Zero 將紀錄 ip 對應 count map 重新初始化
func (l *counterLimit) Zero() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.ip2count = make(map[string]int)
}

// TakeAvailableAndIncr 確認是否可用 (尚未到達限流條件)，並同時累加 1，代表佔一個使用量
func (l *counterLimit) TakeAvailableAndIncr(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if tools.IsEmpty(ip) {
		return false
	}

	isIncr := l.ip2count[ip] < l.maxCount
	if isIncr {
		l.incr(ip)
	}

	return isIncr
}

// incr 計數器累加 1
func (l *counterLimit) incr(ip string) {
	l.ip2count[ip] += 1
}

// GetCount 依照 IP 參數，取得目前累計數量
func (l *counterLimit) GetCount(ip string) string {
	l.mu.Lock()
	defer l.mu.Unlock()

	if tools.IsEmpty(ip) {
		return "0"
	}

	return strconv.Itoa(l.ip2count[ip])
}
