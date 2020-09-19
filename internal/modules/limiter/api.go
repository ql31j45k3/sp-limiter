package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/sp-limiter/configs"
	"net/http"
	"strconv"
)

func RegisterRouter(r *gin.Engine) {
	r.GET("/", apiFunc[configs.ConfigHost.GetMode()])
}

func newLimiter() limiterRouter {
	return limiterRouter{}
}

type limiterRouter struct {
}

func (l *limiterRouter) getCountLimiter(c *gin.Context) {
	clientIP := c.ClientIP()

	if countLimit.IsAvailable(clientIP) {
		countLimit.Increase(clientIP)
		c.String(http.StatusOK, countLimit.GetCount(clientIP))
		return
	}

	c.String(http.StatusOK, "Error")
}

func (l *limiterRouter) getTokenBucket(c *gin.Context) {
	clientIP := c.ClientIP()
	block := false

	if ok, count := tokenBucket.TakeAvailable(clientIP, block); ok {
		c.String(http.StatusOK, strconv.Itoa(int(count)))
		return
	}

	c.String(http.StatusOK, "Error")
}
