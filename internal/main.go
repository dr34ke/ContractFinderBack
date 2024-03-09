package main

import (
	routes "contractfinder/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	//r.PATCH("/users")
	r.Use(cors.Default())

	routes.AuthRoutes(r)

	routes.UserRoutes(r)
	routes.OfferRoutes(r)

	r.Run(":8080")
}
