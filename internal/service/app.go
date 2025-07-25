package service

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	DB             *gorm.DB
	Rdb            *redis.Client
	BackendForRust *BackendForRust
	FrontendForWeb *BackendForWeb
}

const (
	endPoint    = "ws://175.27.250.177:8080"
	dnsEndPoint = "root:123456@tcp(172.30.210.158:3306)/miniRustpbxgo?charset=utf8mb4&parseTime=True&loc=Local"
	callType    = "webrtc"
)

func NewApp() *App {
	app := new(App)
	connDB(app)
	connRdb(app)
	app.BackendForRust = NewBackendForRust(endPoint)
	app.FrontendForWeb = NewBackendForWebByNoParam(app.DB)
	go app.BackendForRust.Connect(callType, app.FrontendForWeb)
	return app
}

func connDB(app *App) {
	dns := dnsEndPoint
	db, err := gorm.Open(mysql.Open(dns))
	if err != nil {
		panic(err)
	}
	mysqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	mysqlDB.SetMaxIdleConns(10)
	mysqlDB.SetMaxOpenConns(100)
	app.DB = db
}

func connRdb(app *App) {
	// redis-cli
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	app.Rdb = rdb
}
