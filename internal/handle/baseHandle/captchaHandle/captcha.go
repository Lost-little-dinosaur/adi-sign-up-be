package captchaHandle

import (
	err2 "adi-sign-up-be/internal/dto/err"
	"adi-sign-up-be/internal/middleware"
	"adi-sign-up-be/pkg/utils/captcha"
	"github.com/gin-gonic/gin"
)

func HandleGetCaptcha(c *gin.Context) {
	cp, err := captcha.GenerateCaptcha()
	if err != nil {
		middleware.Fail(c, err2.InternalErr)
		return
	}
	middleware.Success(c, *cp)
}
