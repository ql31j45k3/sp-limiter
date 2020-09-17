package limiter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/sp-limiter/configs"
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
