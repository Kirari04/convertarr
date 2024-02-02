package app

import (
	"encoder/m"

	"gorm.io/gorm"
)

var Hostname string = "0.0.0.0"
var Port string = "8080"
var Name string = "Convertarr"
var DB *gorm.DB
var Setting *m.SettingValue
