package v1

import (
	"adi-sign-up-be/config"
	"adi-sign-up-be/internal/middleware"
	"adi-sign-up-be/internal/router/v1/baseServiceRouter"
	"adi-sign-up-be/internal/router/v1/signUpRouter"
	"fmt"
	"github.com/gin-gonic/gin"
)

func MainRouter(e *gin.Engine) {
	e.Any("", func(c *gin.Context) {
		data := struct {
			UA         string
			Host       string
			Method     string
			Proto      string
			RemoteAddr string
			Message    string
		}{
			UA:         c.Request.Header.Get("User-Agent"),
			Host:       c.Request.Host,
			Method:     c.Request.Method,
			Proto:      c.Request.Proto,
			RemoteAddr: c.Request.RemoteAddr,
			Message:    fmt.Sprintf("Welcome to %s, version %s.", config.GetConfig().ProgramName, config.GetConfig().VERSION),
		}
		middleware.Success(c, data)
	})
	baseServiceRouter.InitBaseServiceRouter(e)
	signUpRouter.InitSignUpRouter(e)
}
