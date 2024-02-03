package helper

import (
	"encoder/t"

	"github.com/labstack/echo/v4"
)

// Templ CTX
func TCtx(c echo.Context) t.TemplCtx {
	isAuth, _ := c.Get("IsAuth").(bool)
	return t.TemplCtx{
		Error:  nil,
		IsAuth: isAuth,
	}
}

// Templ CTX with error
func TCtxWError(c echo.Context, err error) t.TemplCtx {
	isAuth, _ := c.Get("IsAuth").(bool)
	return t.TemplCtx{
		Error:  err,
		IsAuth: isAuth,
	}
}
