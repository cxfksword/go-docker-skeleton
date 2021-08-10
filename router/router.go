package router

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Create(f *embed.FS) *gin.Engine {
	r := gin.Default()

	// 配置路由
	r.StaticFile("/", "view/dist/index.html")
	r.StaticFile("/favicon.ico", "view/dist/favicon.ico")
	r.GET("/static/*file", func(c *gin.Context) {
		file := c.Param("file")
		c.FileFromFS(fmt.Sprintf("view/dist/static/%s", file), http.FS(f))
	})

	return r
}
