package routes

import (
	controller "contractfinder/internal/controllers"
	"contractfinder/internal/middleware"

	"github.com/gin-gonic/gin"
)

func OfferRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("categories", controller.GetCategories())
	incomingRoutes.GET("categories-names", controller.GetCategoriesNames())
	incomingRoutes.GET("category-offers/:id", controller.GetCategoryOffers())
	incomingRoutes.GET("offer/:id", controller.GetOffer())
	incomingRoutes.POST("offer-apply", controller.UserApplication())
	incomingRoutes.POST("offer", controller.AddOffer())
}
