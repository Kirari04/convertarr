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

var JwtSecret string = ""
var FilesToEncode []string
var LastScanNFiles uint64
var LastFileScan *time.Time
var IsFileScanning bool
var CurrentFileToEncode string
var Cache = cache.New(5*time.Minute, 10*time.Minute)

var AwaitForFileCopy string
var AwaitForFileCopyChan = make(chan string)

var PreloadedFiles *t.TPreloadedFiles = &t.TPreloadedFiles{
	List: []*t.PreloadedFile{},
}
