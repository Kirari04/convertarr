package setup

import (
	"encoder/app"
	"encoder/m"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

func Migrate() {
	if app.DB == nil {
		log.Fatal("DB global is nil")
	}

	m := gormigrate.New(
		app.DB, gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			{
				ID: "3",
				Migrate: func(tx *gorm.DB) error {
					type Setting struct {
						ID        uint `gorm:"primarykey"`
						CreatedAt time.Time
						UpdatedAt time.Time
						Value     m.SettingValue
					}

					return tx.Migrator().CreateTable(&Setting{})
				},
				Rollback: func(tx *gorm.DB) error {
					return tx.Migrator().DropTable("settings")
				},
			},
			{
				ID: "4",
				Migrate: func(tx *gorm.DB) error {
					type Folder struct {
						ID        uint `gorm:"primarykey"`
						CreatedAt time.Time
						UpdatedAt time.Time
						Path      string
					}

					return tx.Migrator().CreateTable(&Folder{})
				},
				Rollback: func(tx *gorm.DB) error {
					return tx.Migrator().DropTable("folders")
				},
			},
		},
	)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Info("Migration did run successfully")

}
