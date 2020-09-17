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

	ip := []string{"192.0.2.1:1235", "192.0.2.2:1235"}
	for i := 0; i < len(ip); i++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup, t *testing.T, r *gin.Engine, url, remoteAddr string) {
			t.Log(fmt.Sprintf("remoteAddr : %s", remoteAddr))
			defer wg.Done()

			testLimiterUnblocked(t, r, url, remoteAddr)
			testLimiterBlock(t, r, url, remoteAddr)

			// 睡眠模擬被阻擋一段時間
			t.Log(fmt.Sprintf("限流 %s 模擬被阻擋一段時間", configs.ConfigHost.GetInterval()))
			ticker := time.NewTicker(configs.ConfigHost.GetInterval())

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
	for i := 0; i < configs.ConfigHost.GetMaxCount(); i++ {
		httpStatus, body, err := httptestRequest(r, http.MethodGet, url, nil, remoteAddr)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, http.StatusOK, httpStatus)
		response := i + 1
		assert.Equal(t, string(body), strconv.Itoa(response))
	}
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
