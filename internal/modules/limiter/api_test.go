package limiter

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ql31j45k3/sp-limiter/configs"

	"github.com/stretchr/testify/assert"
)

func start() *gin.Engine {
	path, err2 := os.Getwd()
	if err2 != nil {
		panic(err2)
	}
	// 測試執行起點位置不一樣，先手動調整取得路徑，才可正常取得 config.yaml 設定檔
	path = path[0:strings.Index(path, "sp-limiter")] + "sp-limiter"
	configs.Start(path)

	Start()

	r := gin.Default()
	RegisterRouter(r)

	return r
}

// httptestRequest 根據特定請求 URL 和參數 param
func httptestRequest(r *gin.Engine, method, uri string, reader io.Reader, remoteAddr string) (int, []byte, error) {
	req := httptest.NewRequest(method, uri, reader)
	// 模擬不同 IP/客戶端
	req.RemoteAddr = remoteAddr

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return 0, nil, err
	}

	return w.Code, body, nil
}

func TestRegisterRouter(t *testing.T) {
	r := start()

	url, err := url.Parse("/")
	if err != nil {
		t.Error(err)
		return
	}

	var wg sync.WaitGroup

	// 模擬兩個以上 Client，每個 IP 有總數量上限
	ip := []string{"192.0.2.1:1235", "192.0.2.2:1235"}

	// 測試資料 10 筆
	//ip := []string{"192.0.2.1:1235", "192.0.2.2:1235", "192.0.2.3:1235", "192.0.2.4:1235", "192.0.2.5:1235",
	//	"192.0.2.6:1235", "192.0.2.7:1235", "192.0.2.8:1235", "192.0.2.9:1235", "192.0.2.10:1235"}
	for i := 0; i < len(ip); i++ {
		wg.Add(1)

		// 模擬請求並發送出，模擬同時多 IP/客戶端
		go func(wg *sync.WaitGroup, t *testing.T, r *gin.Engine, url, remoteAddr string) {
			defer wg.Done()

			testLimiterUnblocked(t, r, url, remoteAddr)
			testLimiterBlock(t, r, url, remoteAddr)

			// 睡眠模擬被阻擋一段時間
			t.Log(fmt.Sprintf("限流 %s 模擬被阻擋一段時間", configs.ConfigHost.GetInterval()))
			// 原本阻擋再加 1s 還是可能有時間差問題
			ticker := time.NewTicker(configs.ConfigHost.GetInterval() + 1)

			<-ticker.C
			// 超過限流限制時間，重新發起請求
			testLimiterUnblocked(t, r, url, remoteAddr)
			testLimiterBlock(t, r, url, remoteAddr)

		}(&wg, t, r, url.String(), ip[i])
	}

	wg.Wait()
}

// testLimiterUnblocked 限流尚未到達上限，正常取得 requests 數量
func testLimiterUnblocked(t *testing.T, r *gin.Engine, url, remoteAddr string) {
	var wg2 sync.WaitGroup

	for i := 0; i < configs.ConfigHost.GetMaxCount(); i++ {
		wg2.Add(1)

		// 模擬請求並發送出，模擬同個 IP/客戶端執行多次 API 請求
		go func(wg2 *sync.WaitGroup, t *testing.T, r *gin.Engine, url, remoteAddr string) {
			defer wg2.Done()

			httpStatus, body, err := httptestRequest(r, http.MethodGet, url, nil, remoteAddr)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, http.StatusOK, httpStatus)

			// 用條件判斷，只要目前數量未超過限制數量，代表正確
			assert.Condition(t, func() bool {
				responseCount, err := strconv.Atoi(string(body))
				if err != nil {
					t.Error(err)
				}
				
				t.Log(fmt.Sprintf("remoteAddr %s, maxCount %d, responseCount %d, maxCount > responseCount %t",
					remoteAddr, configs.ConfigHost.GetMaxCount(), responseCount, responseCount > configs.ConfigHost.GetMaxCount()))

				if responseCount > configs.ConfigHost.GetMaxCount() {
					return false
				}

				return true
			}, fmt.Sprintf("限流未成功阻擋，超過限制 %d 數量的上限", configs.ConfigHost.GetMaxCount()))
		}(&wg2, t, r, url, remoteAddr)
	}

	wg2.Wait()
}

// testLimiterBlock 限流到達上限，阻止正常請求
func testLimiterBlock(t *testing.T, r *gin.Engine, url, remoteAddr string) {
	httpStatus, body, err := httptestRequest(r, http.MethodGet, url, nil, remoteAddr)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, http.StatusOK, httpStatus)
	assert.Equal(t, string(body), "Error")
}
