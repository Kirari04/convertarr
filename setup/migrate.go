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
			{
				ID: "5",
				Migrate: func(tx *gorm.DB) error {
					type History struct {
						ID        uint `gorm:"primarykey"`
						CreatedAt time.Time
						UpdatedAt time.Time
						OldPath   string
						NewPath   string
						OldSize   uint64
						NewSize   uint64
						Error     string `gorm:"size:10000"`
						Status    string // encoding | failed | finished | copy
					}

					return tx.Migrator().CreateTable(&History{})
				},
				Rollback: func(tx *gorm.DB) error {
					return tx.Migrator().DropTable("histories")
				},
			},
		},
	)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Info("Migration did run successfully")

}
