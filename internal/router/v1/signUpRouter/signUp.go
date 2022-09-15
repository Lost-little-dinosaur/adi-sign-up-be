package signUpRouter

import (
	"adi-sign-up-be/internal/handle/signUpHandle"
	"github.com/gin-gonic/gin"
)

func InitSignUpRouter(e *gin.Engine) {
	baseGroup := e.Group("/signup")
	{
		baseGroup.POST("/signUp", signUpHandle.HandleAddSignUp)
		baseGroup.GET("/getAllSignUp", signUpHandle.HandleGetAllSignUp)
	}
}
