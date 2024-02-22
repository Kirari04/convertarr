package app

import (
	"encoder/m"
	"encoder/t"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/patrickmn/go-cache"
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
var LastScanNFiles int
var LastFileScan *time.Time
var CurrentFileToEncode string
var Cache = cache.New(5*time.Minute, 10*time.Minute)
