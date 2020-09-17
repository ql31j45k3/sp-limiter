package limiter

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRouter(r *gin.Engine) {
	limiterRouter := newLimiter()

	routerGroup := r.Group("/")
	routerGroup.GET("", limiterRouter.get)
}

func newLimiter() limiterRouter {
	return limiterRouter{}
}

type limiterRouter struct {
}

func (l *limiterRouter) get(c *gin.Context) {
	clientIP := c.ClientIP()
	c.String(http.StatusOK, "hihi " + clientIP)
}
