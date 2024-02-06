package server

import (
	"context"
	"encoder/app"
	"encoder/handler"
	"encoder/mware"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Serve() {
	var Address = fmt.Sprintf("%s:%s", app.Hostname, app.Port)

	// server config
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(log.INFO)

	e.Use(mware.HasBeenSetupRedirect)
	e.Use(mware.AuthEnabledRedirect)

	// routes
	e.GET("/", handler.GetIndex)
	e.GET("/setup", handler.GetSetup)
	e.POST("/setup", handler.PostSetup)
	e.POST("/scanner", handler.PostScanner)
	e.GET("/history", handler.GetHistory)
	e.GET("/history/table", handler.GetHistoryTable)
	e.GET("/setting", handler.GetSetting)
	e.POST("/setting", handler.PostSetting)
	e.GET("/setting/user", handler.GetSettingUser)
	e.POST("/setting/user", handler.PostSettingUser)
	e.GET("/setting/folder", handler.GetSettingFolder)
	e.POST("/setting/folder", handler.PostSettingFolder)
	e.POST("/setting/folder/delete", handler.DeleteSettingFolder)
	e.GET("/login", handler.GetLogin)
	e.POST("/login", handler.PostLogin)
	e.POST("/logout", handler.PostLogout)

	// start & shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(Address); err != nil {
			if err != http.ErrServerClosed {
				e.Logger.Info("server crashed", err)
			} else {
				e.Logger.Info("shutting down the server")
			}

		}
	}()
	log.Printf("Server started on address: http://%s", Address)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
