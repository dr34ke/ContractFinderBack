package controllers

import (
	"contractfinder/internal/database"
	helper "contractfinder/internal/helpers"
	"contractfinder/internal/models"
	"log"
	"net/http"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SingUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

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

		countEmail, err := database.Count("user", bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		countPhone, err := database.Count("user", bson.M{"phone": user.Phone})
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

		result, err := database.Insert("user", user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		foundUser, err := database.GetOne[models.User](database.DBinstance(), "user", bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusBadRequest, "Nie znaleziono użytkownika")
			log.Print(err)
			return
		}

		isPasswordCorrect := foundUser.CheckPasswordHash(user.Password)
		if !isPasswordCorrect {
			c.JSON(http.StatusBadRequest, "Błędne hasło")
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser)

		foundUser.Token = token
		foundUser.RefreshToken = refreshToken
		foundUser.TimeStamp.Login()
		err = foundUser.UpdateTokens()

		if err != nil {
			c.JSON(http.StatusBadRequest, "Server Error")
			return
		}

		c.JSON(http.StatusOK, foundUser.ReturnSimplified())
	}

}

func GetUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		foundUser, err := database.GetOne[models.User](database.DBinstance(), "user", bson.M{"_id": id})
		if err != nil {
			c.JSON(http.StatusBadRequest, "Nie znaleziono użytkownika")
			log.Print(err)
			return
		}

		c.JSON(http.StatusOK, foundUser.UserProfile)
	}
}
func GetUserPreference() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		foundUser, err := database.GetOne[models.User](database.DBinstance(), "user", bson.M{"_id": id})
		if err != nil {
			c.JSON(http.StatusBadRequest, "Nie znaleziono użytkownika")
			log.Print(err)
			return
		}

		c.JSON(http.StatusOK, foundUser.UserPreference)
	}
}

func UpdateUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		response, err := database.Update("user", filter, update)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func UpdateUserPreference() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		response, err := database.Update("user", filter, update)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
func GetUserPublicProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		pipeline := mongo.Pipeline{
			bson.D{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "UserRatings"},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "idEvaluatedUser"},
				{Key: "as", Value: "ratingModel"},
			}}},
			bson.D{{Key: "$match", Value: bson.D{
				{Key: "_id", Value: bson.M{"$eq": id}},
			}}},
			bson.D{{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "firstName", Value: 1},
				{Key: "lastName", Value: 1},
				{Key: "email", Value: 1},
				{Key: "phone", Value: 1},
				{Key: "userProfile", Value: 1},
				{Key: "userPreference", Value: 1},
				{Key: "rating", Value: bson.M{"$avg": "$ratingModel.rating"}},
			}}},
		}

		results, err := database.Aggregate[models.User]("user", pipeline)

		if err != nil {
			log.Print(err)
			c.JSON(http.StatusBadRequest, "Nie znaleziono ofert"+err.Error())
			return
		}

		c.JSON(http.StatusOK, results[0].ReturnUserDTO())
	}
}

func GetUserRatings() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		userRatings, err := database.GetMany[models.UserRating]("UserRatings", bson.M{"idEvaluatedUser": id})

		if err != nil {
			log.Print(err)
			c.JSON(http.StatusBadRequest, "")
			return
		}
		c.JSON(http.StatusOK, userRatings)
	}
}
