package main

import (
	routes "contractfinder/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	//r.PATCH("/users")

	routes.AuthRoutes(r)
	routes.UserRoutes(r)

	r.Run("localhost:8081")
}
