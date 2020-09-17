package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/sp-limiter/configs"
	"net/http"
)

func RegisterRouter(r *gin.Engine) {
	routerGroup := r.Group("/")
	routerGroup.GET("", apiFunc[configs.ConfigHost.GetMode()])
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
		c.String(http.StatusOK, "clientIP : "+clientIP+" Request count: "+countLimit.GetCount(clientIP))
		return
	}

	c.String(http.StatusOK, "Error")
}
