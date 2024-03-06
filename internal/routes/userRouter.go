package routes

import (
	controller "contractfinder/internal/controllers"
	"contractfinder/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("user/update-profile", controller.UpdateUserProfile())
	incomingRoutes.POST("user/update-preference", controller.UpdateUserPreference())
	incomingRoutes.GET("user/get-profile/:id", controller.GetUserProfile())
	incomingRoutes.GET("user/get-preference/:id", controller.GetUserPreference())
}
