package setup

import (
	"encoder/app"
	"encoder/m"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
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
				ID: "1",
				Migrate: func(tx *gorm.DB) error {
					// it's a good pratice to copy the struct inside the function,
					// so side effects are prevented if the original struct changes during the time
					type User struct {
						ID        uuid.UUID `gorm:"type:uuid;primaryKey;uniqueIndex"`
						CreatedAt time.Time
						UpdatedAt time.Time
						Username  string
						PwdHash   string
					}
					return tx.Migrator().CreateTable(&User{})
				},
				Rollback: func(tx *gorm.DB) error {
					return tx.Migrator().DropTable("users")
				},
			},
			{
				ID: "2",
				Migrate: func(tx *gorm.DB) error {
					type UserSession struct {
						ID        uuid.UUID `gorm:"type:uuid;primaryKey;uniqueIndex"`
						CreatedAt time.Time
						UpdatedAt time.Time
						User      m.User
						UserId    uuid.UUID `gorm:"index"`
					}

					return tx.Migrator().CreateTable(&UserSession{})
				},
				Rollback: func(tx *gorm.DB) error {
					return tx.Migrator().DropTable("user_sessions")
				},
			},
			{
				ID: "3",
				Migrate: func(tx *gorm.DB) error {
					type Setting struct {
						ID        uuid.UUID `gorm:"type:uuid;primaryKey;uniqueIndex"`
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
		},
	)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Info("Migration did run successfully")

}
