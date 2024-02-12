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
			c.JSON(http.StatusBadRequest, gin.H{"Validation error": validationErr.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone or email already exists!"})
			return
		}
		user.TimeStamp.Created()
		user.Id = guid.New().String()
		token, refresh_token, _ := helper.GenerateAllTokens(user)
		user.Token = token
		user.Refresh_Token = refresh_token
		user.HashPassword()

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}
func serverError(err error, c *gin.Context) {

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}
		isPasswordCorrect := foundUser.CheckPasswordHash(user.Password)
		if !isPasswordCorrect {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong password"})
			return
		}
		token, refreshToken, _ := helper.GenerateAllTokens(foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		foundUser.Token = token
		foundUser.Refresh_Token = refreshToken
		foundUser.TimeStamp.Login()
		err = foundUser.UpdateTokens(ctx)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Server Error"})
			return
		}

		c.JSON(http.StatusOK, foundUser.ReturnSimplified())
	}

}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": "ok"})
		return
	}
}
