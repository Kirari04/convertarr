package app

import (
	"encoder/m"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var Hostname string = "0.0.0.0"
var Port string = "8080"
var Name string = "Convertarr"
var DB *gorm.DB
var Setting *m.SettingValue
var Validate = validator.New(validator.WithRequiredStructEnabled())
var TemporaryDb bool
