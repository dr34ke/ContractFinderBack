package models

import (
	"contractfinder/internal/database"
	"contractfinder/internal/dto"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             string          `bson:"_id" json:"id"`
	FirstName      string          `bson:"firstName" json:"firstName" validate:"required,min=3,max=25"`
	LastName       string          `bson:"lastName" json:"lastName" validate:"required,min=3,max=25"`
	Password       string          `bson:"password" json:"password"`
	Email          string          `bson:"email" json:"email" validate:"required,min=5,max=100"`
	Phone          string          `bson:"phone" json:"phone" validate:"required,min=9,max=12"`
	Token          string          `bson:"token" json:"token"`
	FirebaseToken  string          `bson:"firebaseToken" json:"firebaseToken"`
	RefreshToken   string          `bson:"refreshToken" json:"refreshToken"`
	TimeStamp      TimeStamp       `bson:"timeStamp" json:"timeStamp"`
	UserProfile    *UserProfile    `bson:"userProfile" json:"userProfile"`
	UserPreference *UserPreference `bson:"userPreference" json:"userPreference"`
	Rating         *float64        `bson:"rating" json:"rating"`
}

type SimplifiedUser struct {
	Id            string `bson:"_id" json:"id"`
	First_name    string `bson:"first_name" json:"first_name"`
	Last_name     string `bson:"last_name" json:"last_name"`
	Email         string `bson:"email" json:"email"`
	Phone         string `bson:"phone" json:"phone"`
	Token         string `bson:"token" json:"token"`
	Refresh_Token string `bson:"refresh_token" json:"refresh_token"`
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.Password = string(bytes)
	return err
}

func (u User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) UpdateTokens() error {
	filter := bson.M{"_id": u.Id}
	update := bson.M{
		"$set": bson.M{
			"token":                u.Token,
			"refreshToken":         u.RefreshToken,
			"timeStamp.last_login": u.TimeStamp.Last_login,
		},
	}
	_, err := database.Update("user", filter, update)
	return err
}

func (u User) ReturnSimplified() SimplifiedUser {
	return SimplifiedUser{
		Id:            u.Id,
		First_name:    u.FirstName,
		Last_name:     u.LastName,
		Email:         u.Email,
		Phone:         u.Phone,
		Token:         u.Token,
		Refresh_Token: u.RefreshToken,
	}
}
func (u User) ReturnUserDTO() dto.UserProfileDTO {

	dto := dto.UserProfileDTO{
		Id:         u.Id,
		First_name: u.FirstName,
		Last_name:  u.LastName,
		Rating:     u.Rating,
	}
	if u.UserPreference != nil && u.UserPreference.IsEmailPublic {
		dto.Email = &u.Email
	}
	if u.UserPreference != nil && u.UserPreference.IsPhonePublic {
		dto.Phone = &u.Phone
	}

	if u.UserProfile != nil {
		dto.Description = u.UserProfile.Description
		dto.Image = u.UserProfile.Image
	}

	return dto
}
