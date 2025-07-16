package setup

import (
	"encoder/app"
	"encoder/helper"
	"log"
)

func App() {
	randomString, err := helper.RandomString(32)
	if err != nil {
		log.Fatalf("Failed to generate random string: %v", err)
	}
	app.JwtSecret = randomString
}
