package server

import (
	"adi-sign-up-be/config"
	_ "adi-sign-up-be/internal/corn"
	"adi-sign-up-be/internal/logger"
	"adi-sign-up-be/internal/middleware"
	"adi-sign-up-be/internal/redis"
	v1 "adi-sign-up-be/internal/router/v1"
	"github.com/gin-gonic/gin"
)

var E *gin.Engine

func init() {
	logger.Info.Println("start init gin")
	gin.SetMode(config.GetConfig().MODE)
	E = gin.New()
	E.Use(middleware.GinRequestLog, gin.Recovery(), middleware.Cors(E))
}

func Run() {
	redis.GetRedis()
	v1.MainRouter(E)
	if err := E.Run("0.0.0.0:" + config.GetConfig().PORT); err != nil {
		logger.Error.Fatalln(err)
	}
}
