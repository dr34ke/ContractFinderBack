package helper

import (
	"contractfinder/internal/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SECRET string = os.Getenv("JWT_KEY")

func GenerateAllTokens(user models.User) (singedToken string, signedRefreshToken string, err error) {

	claims := &SignedDetails{
		Email:      user.Email,
		First_name: user.FirstName,
		Last_name:  user.LastName,
		Uid:        user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET))
	refreshToken, err2 := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET))

	if err1 != nil {
		err = err1
	}
	if err2 != nil {
		err = err2
	}

	return token, refreshToken, err
}

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.StandardClaims
}
