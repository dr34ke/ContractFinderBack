package controllers

import (
	"context"
	"contractfinder/internal/database"
	"contractfinder/internal/models"
	"log"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoriesCollection *mongo.Collection = database.OpenConnection(database.DBinstance(), "WorkCategory")

func GetCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		cur, err := categoriesCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.JSON(http.StatusBadRequest, "Nie znaleziono kategorii")
			return
		}

		var results []models.WorkCategory
		if err = cur.All(ctx, &results); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, results)
	}
}
