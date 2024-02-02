package setup

import (
	"encoder/app"

	"github.com/glebarez/sqlite"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

func Db() {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to open database", err)
	}
	app.DB = db
}
