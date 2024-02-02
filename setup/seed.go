package setup

import (
	"encoder/app"
	"encoder/m"

	"github.com/labstack/gommon/log"
)

func Seed() {
	// check if settings table is filled
	if app.Setting != nil {
		log.Fatal("App Settings are already instantiated")
	}
	var setting m.Setting
	if err := app.DB.FirstOrCreate(&setting).Error; err != nil {
		log.Fatal("Failed to instantiated db settings")
	}
	app.Setting = &setting.Value
}
