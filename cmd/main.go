package main

import (
	"github.com/gin-gonic/gin"
	"miniRustpbxgo/internal/api"
	"miniRustpbxgo/internal/service"
)

func main() {
	app := service.NewApp()
	// backendForWeb和backendForRust两者要做好初始化
	api.Routers(gin.Default(), app)
	select {}
}
