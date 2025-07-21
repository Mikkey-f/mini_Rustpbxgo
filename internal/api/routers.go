package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"miniRustpbxgo/internal/service"
)

const (
	AuthPath   = "/in"
	NoAuthPath = "/out"
)

func Routers(router *gin.Engine, backendForWeb *service.BackendForWeb, backendForRust *service.BackendForRust) {
	//auth := router.Group(AuthPath)
	//{
	//
	//}
	noAuth := router.Group(NoAuthPath)
	{
		noAuth.GET("/webrtc/setup", func(c *gin.Context) {
			backendForWeb.HandleWebRtcSetUp(c.Writer, c.Request)
		})
	}

	// 防止阻塞
	go func() {
		// 阻塞进程
		if err := router.Run(":8081"); err != nil {
			log.Fatal(err, "路由建立失败")
		}
	}()

}
