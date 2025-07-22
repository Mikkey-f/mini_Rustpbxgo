package service

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	db *gorm.DB
}

func NewApp() *App {
	app := new(App)
	connDB(app)
	return app
}

func connDB(app *App) {
	dns := "root:123456@tcp(172.30.210.158:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
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
	app.db = db
}
