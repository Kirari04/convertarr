package server

import (
	"context"
	"encoder/app"
	"encoder/views"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Serve() {
	var Address = fmt.Sprintf("%s:%s", app.Hostname, app.Port)
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(log.INFO)
	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, views.Index())
	})
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

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}
