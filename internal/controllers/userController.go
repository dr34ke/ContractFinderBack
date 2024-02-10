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
		serverError(err, c)

		countPhone, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		serverError(err, c)

		if countEmail > 0 || countPhone > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone or email already exists!"})
		}
		user.TimeStamp.Created()
		user.Id = guid.New()
		token, refresh_token, _ := helper.GenerateAllTokens(user)
		user.Token = token
		user.Refresh_Token = refresh_token

		result, err := userCollection.InsertOne(ctx, user)
		serverError(err, c)
		c.JSON(http.StatusCreated, result)
	}
}
func serverError(err error, c *gin.Context) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUser() {

}
