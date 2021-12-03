package route

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InitRoute 路由入口
func InitRoute(r *gin.Engine) {

	// listen 存过监控
	r.Any("/listen", func(c *gin.Context) {
		c.String(200, "Success")
	})

	// metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
