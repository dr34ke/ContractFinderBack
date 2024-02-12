package routes

import (
	controller "contractfinder/internal/controllers"
	"contractfinder/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("user/test", controller.GetUser())
}
