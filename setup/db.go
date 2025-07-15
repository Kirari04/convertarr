package setup

import (
	"encoder/app"
	"os"

	"github.com/labstack/gommon/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Db() {
	var databasePath = "./database/db.sqlite"
	if app.TemporaryDb {
		log.Info("Using temporary database")
		databasePath = "file::memory:?cache=shared"
	} else {
		if err := os.MkdirAll("./database", 0766); err != nil {
			log.Fatal("Failed to create db file", err)
		}
	}
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to open database", err)
	}
	app.DB = db
}
