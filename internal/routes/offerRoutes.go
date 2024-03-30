package routes

import (
	controller "contractfinder/internal/controllers"
	"contractfinder/internal/middleware"

	"github.com/gin-gonic/gin"
)

func OfferRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("categories", controller.GetCategories())
	incomingRoutes.GET("category-offers/:id", controller.GetCategoryOffers())
	incomingRoutes.GET("offer/:id", controller.GetOffer())
}
