package setup

import (
	"encoder/app"

	"github.com/glebarez/sqlite"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

func Db() {
	var databasePath = "./database/db.sqlite"
	if app.TemporaryDb {
		log.Info("Using temporary database")
		databasePath = "file::memory:?cache=shared"
	}
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to open database", err)
	}
	app.DB = db
}
