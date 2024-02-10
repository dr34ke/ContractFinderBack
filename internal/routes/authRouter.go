package routes

import (
	controller "contractfinder/internal/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("user/singup", controller.SingUp())
	incomingRoutes.POST("user/login", controller.Login())
}
