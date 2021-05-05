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
	r.GET("/", func(c *gin.Context) {
		file, _ := f.ReadFile("view/dist/index.html")
		c.Data(
			http.StatusOK,
			"text/html; charset=utf-8",
			file,
		)
	})

	r.GET("/public/:name", func(c *gin.Context) {
		name := c.Param("name")
		fmt.Println(name)
		c.FileFromFS(fmt.Sprintf("view/dist/%s", name), http.FS(f))
	})

	return r
}
