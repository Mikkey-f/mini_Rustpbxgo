package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"miniRustpbxgo/internal/api/filter"
	"miniRustpbxgo/internal/service"
)

const (
	AuthPath   = "/in"
	NoAuthPath = "/out"
)

func Routers(router *gin.Engine, app *service.App) {
	//auth := router.Group(AuthPath)
	//{
	//
	//}
	authFilter := filter.NewSessionAuth()
	auth := router.Group(AuthPath).Use(authFilter.Auth)
	noAuth := router.Group(NoAuthPath)
	{

		noAuth.POST("/user/register", app.Register)
		noAuth.GET("/user/login", app.Login)
		noAuth.GET("/webrtc/setup", func(c *gin.Context) {
			app.FrontendForWeb.HandleWebRtcSetUp(c.Writer, c.Request, app.BackendForRust)
		})
	}
	{

		auth.POST("/create/robotKey", app.CreateRobotKey)
		auth.GET("/list/robotKey", app.RobotKeyList)
		auth.POST("/create/robot", app.CreateRobot)
		auth.GET("/list/robot", app.RobotList)
		auth.PUT("/update/robot", app.UpdateRobot)
	}

	// 防止阻塞
	go func() {
		// 阻塞进程
		if err := router.Run(":8081"); err != nil {
			log.Fatal(err, "路由建立失败")
		}
	}()

}
