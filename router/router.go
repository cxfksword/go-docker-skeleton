package router

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Create(f *embed.FS) *gin.Engine {
	r := gin.Default()

	// 配置路由
	r.GET("/", func(c *gin.Context) {
		file, _ := f.ReadFile("view/index.html")
		c.Data(
			http.StatusOK,
			"text/html; charset=utf-8",
			file,
		)
	})

	r.StaticFS("/public", http.FS(f))

	return r
}
