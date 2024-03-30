package controllers

import (
	"context"
	"contractfinder/internal/database"
	"contractfinder/internal/models"
	"log"

	"net/http"
	"time"

	"strconv"

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

		longitude, _ := strconv.ParseFloat(c.Query("longitude"), 64)
		latitude, _ := strconv.ParseFloat(c.Query("latitude"), 64)
		distance, _ := strconv.ParseFloat(c.Query("distance"), 64)

		distance = distance / 6378.1

		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}

		lookupStage := bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "WorkOffers"},
				{Key: "localField", Value: "_id"},
				{Key: "pipeline", Value: bson.A{
					bson.M{
						"$match": bson.D{
							{Key: "$or", Value: bson.A{
								bson.M{
									"coordinates": bson.D{
										{Key: "$geoWithin", Value: bson.D{
											{Key: "$centerSphere", Value: bson.A{
												bson.A{latitude, longitude},
												distance,
											}}}}}},
								bson.M{"onSite": false},
							}}}}}},
				{Key: "foreignField", Value: "categoryId"},
				{Key: "as", Value: "offers"},
			}}}

		projectStage := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "name", Value: 1},
			{Key: "offersCount", Value: bson.M{"$size": "$offers"}},
		}}}

		// Aggregate pipeline
		pipeline := mongo.Pipeline{matchStage, lookupStage, projectStage}

		cur, err := categoriesCollection.Aggregate(context.Background(), pipeline)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Nie znaleziono kategorii"+err.Error())
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

		longitude, _ := strconv.ParseFloat(c.Query("longitude"), 64)
		latitude, _ := strconv.ParseFloat(c.Query("latitude"), 64)
		distance, _ := strconv.ParseFloat(c.Query("distance"), 64)

		pipeline := mongo.Pipeline{
			bson.D{{Key: "$geoNear", Value: bson.D{
				{Key: "near", Value: bson.D{
					{Key: "type", Value: "Point"},
					{Key: "coordinates", Value: bson.A{latitude, longitude}},
				}},
				{Key: "distanceField", Value: "distance"},
				{Key: "maxDistance", Value: distance * 1000},
				{Key: "spherical", Value: true},
			}}},
			bson.D{{Key: "$match", Value: bson.D{
				{Key: "categoryId", Value: bson.M{"$eq": id}},
			}}},
			bson.D{{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "userId", Value: 1},
				{Key: "categoryId", Value: 1},
				{Key: "title", Value: 1},
				{Key: "description", Value: 1},
				{Key: "sugestedSalary", Value: 1},
				{Key: "isSalaryPerHour", Value: 1},
				{Key: "isRepetetive", Value: 1},
				{Key: "isFromWorker", Value: 1},
				{Key: "coordinates", Value: 1},
				{Key: "onSite", Value: 1},
				{Key: "distanceInKm", Value: bson.M{"$divide": bson.A{"$distance", 1000}}},
			}}},
			bson.D{{Key: "$unionWith", Value: bson.D{
				{Key: "coll", Value: "WorkOffers"},
				{Key: "pipeline", Value: bson.A{
					bson.M{"$match": bson.M{
						"$and": bson.A{
							bson.M{"categoryId": bson.M{"$eq": id}},
							bson.M{"onSite": bson.M{"$eq": false}},
						},
					}},
					bson.D{{Key: "$project", Value: bson.D{
						{Key: "_id", Value: 1},
						{Key: "userId", Value: 1},
						{Key: "categoryId", Value: 1},
						{Key: "title", Value: 1},
						{Key: "description", Value: 1},
						{Key: "sugestedSalary", Value: 1},
						{Key: "isSalaryPerHour", Value: 1},
						{Key: "isRepetetive", Value: 1},
						{Key: "isFromWorker", Value: 1},
						{Key: "coordinates", Value: 1},
						{Key: "onSite", Value: 1},
					}}},
				}},
			}}},
			bson.D{{Key: "$sort", Value: bson.D{
				{Key: "distanceInKm", Value: 1},
			}}},
		}

		// Run the aggregation pipeline
		cur, err := offersCollection.Aggregate(context.Background(), pipeline)

		if err != nil {
			c.JSON(http.StatusBadRequest, "Nie znaleziono ofert"+err.Error())
			return
		}

		var results []models.WorkOffer
		if err = cur.All(ctx, &results); err != nil {
			log.Print(err)
		}
		c.JSON(http.StatusOK, results)
	}
}
func GetOffer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		id := c.Param("id")

		longitude, _ := strconv.ParseFloat(c.Query("longitude"), 64)
		latitude, _ := strconv.ParseFloat(c.Query("latitude"), 64)

		pipeline := mongo.Pipeline{
			bson.D{{Key: "$geoNear", Value: bson.D{
				{Key: "near", Value: bson.D{
					{Key: "type", Value: "Point"},
					{Key: "coordinates", Value: bson.A{latitude, longitude}},
				}},
				{Key: "distanceField", Value: "distance"},
				{Key: "spherical", Value: true},
			}}},
			bson.D{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "UserApplications"},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "offerId"},
				{Key: "as", Value: "usersApplications"},
			}}},
			bson.D{{Key: "$match", Value: bson.D{
				{Key: "_id", Value: bson.M{"$eq": id}},
			}}},
			bson.D{{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "userId", Value: 1},
				{Key: "categoryId", Value: 1},
				{Key: "title", Value: 1},
				{Key: "description", Value: 1},
				{Key: "sugestedSalary", Value: 1},
				{Key: "isSalaryPerHour", Value: 1},
				{Key: "isRepetetive", Value: 1},
				{Key: "isFromWorker", Value: 1},
				{Key: "coordinates", Value: 1},
				{Key: "onSite", Value: 1},
				{Key: "distanceInKm", Value: bson.M{"$divide": bson.A{"$distance", 1000}}},
				{Key: "usersApplications", Value: 1},
			}}},
			bson.D{{Key: "$unionWith", Value: bson.D{
				{Key: "coll", Value: "WorkOffers"},
				{Key: "pipeline", Value: bson.A{
					bson.M{"$match": bson.M{
						"$and": bson.A{
							bson.M{"_id": bson.M{"$eq": id}},
							bson.M{"onSite": bson.M{"$eq": false}},
						},
					}},
					bson.D{{Key: "$project", Value: bson.D{
						{Key: "_id", Value: 1},
						{Key: "userId", Value: 1},
						{Key: "categoryId", Value: 1},
						{Key: "title", Value: 1},
						{Key: "description", Value: 1},
						{Key: "sugestedSalary", Value: 1},
						{Key: "isSalaryPerHour", Value: 1},
						{Key: "isRepetetive", Value: 1},
						{Key: "isFromWorker", Value: 1},
						{Key: "coordinates", Value: 1},
						{Key: "onSite", Value: 1},
					}}},
				}},
			}}},
		}

		// Run the aggregation pipeline
		cur, err := offersCollection.Aggregate(context.Background(), pipeline)

		if err != nil {
			c.JSON(http.StatusBadRequest, "Nie znaleziono ofert"+err.Error())
			return
		}

		var results []models.WorkOffer
		if err = cur.All(ctx, &results); err != nil {
			log.Print(err)
		}
		c.JSON(http.StatusOK, results[0])
	}
}
