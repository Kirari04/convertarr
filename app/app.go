package app

import (
	"encoder/m"
	"encoder/t"
	"time"

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

var ResourcesHistory t.Resources
var MaxResourcesHistory = 1000
var ResourcesInterval = time.Second * 3
var ResourcesDeleteInterval = time.Minute * 1

var JwtSecret string = "secret"
var FilesToEncode []string
var CurrentFileToEncode string
