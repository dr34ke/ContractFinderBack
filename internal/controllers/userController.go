package controllers

import (
	"context"
	"contractfinder/internal/database"
	helper "contractfinder/internal/helpers"
	"contractfinder/internal/models"

	"net/http"
	"time"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenConnection(database.DBinstance(), "user")
var validate = validator.New()

func SingUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validationErr := validate.Struct(user); validationErr != nil {
			errs := validationErr.(validator.ValidationErrors)
			var response string
			for _, element := range errs {
				if element.Tag() == "min" {
					response += element.Field() + " - za mało znaków, "
				} else {
					response += element.Field() + " - za dużo znaków, "
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": response})
			return
		}

		countEmail, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		countPhone, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if countEmail > 0 || countPhone > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Telefon lub email już istnieje!"})
			return
		}

		user.TimeStamp.Created()
		user.UserPreference.TimeStamp.Created()
		user.UserProfile.TimeStamp.Created()

		user.Id = guid.New().String()
		token, refresh_token, _ := helper.GenerateAllTokens(user)
		user.Token = token
		user.RefreshToken = refresh_token
		user.HashPassword()

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Nie znaleziono użytkownika")
			return
		}
		isPasswordCorrect := foundUser.CheckPasswordHash(user.Password)
		if !isPasswordCorrect {
			c.JSON(http.StatusBadRequest, "Błędne hasło")
			return
		}
		token, refreshToken, _ := helper.GenerateAllTokens(foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		foundUser.Token = token
		foundUser.RefreshToken = refreshToken
		foundUser.TimeStamp.Login()
		err = foundUser.UpdateTokens(ctx)

		if err != nil {
			c.JSON(http.StatusBadRequest, "Server Error")
			return
		}

		c.JSON(http.StatusOK, foundUser.ReturnSimplified())
	}

}

func GetUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		id := c.Param("id")

		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Nie znaleziono użytkownika")
			return
		}

		c.JSON(http.StatusOK, foundUser.UserProfile)
	}
}
func GetUserPreference() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		id := c.Param("id")

		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Nie znaleziono użytkownika")
			return
		}

		c.JSON(http.StatusOK, foundUser.UserPreference)
	}
}

func UpdateUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		id, _ := c.Get("uuid")

		var userProfile models.UserProfile
		if err := c.BindJSON(&userProfile); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		userProfile.TimeStamp.Updated()

		filter := bson.M{"_id": id}
		update := bson.M{
			"$set": bson.M{
				"userProfile": userProfile,
			},
		}
		response, err := userCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func UpdateUserPreference() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		id, _ := c.Get("uuid")

		var userPreference models.UserPreference
		if err := c.BindJSON(&userPreference); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		filter := bson.M{"_id": id}
		update := bson.M{
			"$set": bson.M{
				"userPreference": userPreference,
			},
		}
		response, err := userCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
