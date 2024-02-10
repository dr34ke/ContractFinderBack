package controllers

import (
	"contractfinder/internal/database"

	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenConnection(database.DBinstance(), "user")

func SingUp() {

}

func Login() {

}

func GetUser() {

}
