package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id         string `bson:"_id" json:"id"`
	First_name string `bson:"first_name" json:"first_name" validate:"required min=3, max = 25"`
	Last_name  string `bson:"last_name" json:"last_name" validate:"required min=3, max = 25"`
	Password   string `bson:"password" json:"password"`
	Email      string `bson:"email" json:"email" validate:"required min=5, max = 100"`
	Phone      string `bson:"phone" json:"phone" validate:"required min=9, max = 12"`
	TimeStamp  TimeStamp
}

func (u User) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
