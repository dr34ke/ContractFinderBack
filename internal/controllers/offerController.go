package controllers

import (
	"context"
	"contractfinder/internal/database"
	"contractfinder/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"

	helper "contractfinder/internal/helpers"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

func GetCategoriesNames() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		results, err := database.GetMany[models.WorkCategory]("WorkCategory", bson.M{})

		if err != nil {
			c.JSON(http.StatusBadRequest, "Nie znaleziono kategorii"+err.Error())
			return
		}
		c.JSON(http.StatusOK, results)
	}
}

func GetCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		results, err := database.Aggregate[models.WorkCategory]("WorkCategory", pipeline)

		if err != nil {
			log.Print(err)
			c.JSON(http.StatusBadRequest, "Nie znaleziono kategorii"+err.Error())
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

func GetCategoryOffers() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		results, err := database.Aggregate[models.WorkOffer]("WorkOffers", pipeline)

		if err != nil {
			log.Print(err)
			c.JSON(http.StatusBadRequest, "Nie znaleziono ofert"+err.Error())
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

func GetOffer() gin.HandlerFunc {
	return func(c *gin.Context) {
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
					bson.D{{Key: "$lookup", Value: bson.D{
						{Key: "from", Value: "UserApplications"},
						{Key: "localField", Value: "_id"},
						{Key: "foreignField", Value: "offerId"},
						{Key: "as", Value: "usersApplications"},
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
						{Key: "usersApplications", Value: 1},
					}}},
				}},
			}}},
		}
		results, err := database.Aggregate[models.WorkOffer]("WorkOffers", pipeline)

		if err != nil {
			log.Print(err)
			c.JSON(http.StatusBadRequest, "Nie znaleziono oferty"+err.Error())
			return
		}
		log.Print(pipeline)

		c.JSON(http.StatusOK, results[0])

	}
}

func UserApplication() gin.HandlerFunc {
	return func(c *gin.Context) {
		var application models.UserApplication

		if err := c.BindJSON(&application); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := validate.Struct(application)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			if validationErrors != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
				return
			}
		}

		application.Id = guid.New().String()
		result, err := database.Insert("UserApplications", application)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		offer, err := database.GetOne[models.WorkOffer](database.DBinstance(), "user", bson.M{"_id": application.OfferId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		user, err := database.GetOne[models.User](database.DBinstance(), "user", bson.M{"_id": offer.UserId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		token := user.FirebaseToken

		if err := helper.SendPushNotification(token, "Nowa aplikacja na twoją ofertę", "Ktoś dodał aplikację do oferty którą zamieściłeś. Sprawdź!"); err != nil {
			log.Fatalf("Error sending push notification: %v", err)
		}

		c.JSON(http.StatusCreated, result)
	}
}

func AddOffer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var offer models.WorkOffer

		if err := c.BindJSON(&offer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := validate.Struct(offer)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			if validationErrors != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
				return
			}
		}

		offer.Id = guid.New().String()
		offer.TimeStamp.Created()
		result, err := database.Insert("WorkOffers", offer)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}
