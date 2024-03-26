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
var offersCollection *mongo.Collection = database.OpenConnection(database.DBinstance(), "WorkOffers")

func GetCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		matchStage := bson.D{{"$match", bson.D{}}}

		lookupStage := bson.D{{"$lookup", bson.D{
			{"from", "WorkOffers"},
			{"localField", "_id"},
			{"foreignField", "categoryId"},
			{"as", "offers"},
		}}}

		projectStage := bson.D{{"$project", bson.D{
			{"_id", 1},
			{"name", 1},
			{"offersCount", bson.D{{"$size", "$offers"}}},
		}}}

		// Aggregate pipeline
		pipeline := mongo.Pipeline{matchStage, lookupStage, projectStage}

		cur, err := categoriesCollection.Aggregate(context.Background(), pipeline)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Nie znaleziono kategorii")
			return
		}

		var results []models.WorkCategory
		if err = cur.All(ctx, &results); err != nil {
			log.Print(err)
		}

		c.JSON(http.StatusOK, results)
	}
}

func GetCategoryOffers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		id := c.Param("id")

		cur, err := offersCollection.Find(ctx, bson.D{{"categoryId", id}})

		if err != nil {
			log.Print(err)
			c.JSON(http.StatusBadRequest, "Nie znaleziono ofert")
			return
		}

		var results []models.WorkOffer
		if err = cur.All(ctx, &results); err != nil {
			log.Print(err)
		}
		c.JSON(http.StatusOK, results)
	}
}
