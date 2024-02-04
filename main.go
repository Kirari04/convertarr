package main

import (
	"encoder/app"
	"encoder/server"
	"encoder/setup"
	"fmt"

	"github.com/thatisuday/commando"
)

func main() {
	commando.
		SetExecutableName("Convertarr").
		SetVersion("v1.0.0").
		SetDescription("This CLI tool encodes all mkv files.")

	commando.
		Register(nil).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			fmt.Println("use help command")
		})

	commando.
		Register("serve").
		SetShortDescription("start webserver on :8080").
		AddFlag("dev", "use temporary database", commando.Bool, false).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			dev, _ := flags["dev"].GetBool()
			if dev {
				app.TemporaryDb = true
			}
			setup.Setup()
			server.Serve()
		})

	commando.Parse(nil)
}
