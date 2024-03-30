package dto

type UserProfileDTO struct {
	Id          string   `bson:"_id" json:"id"`
	First_name  string   `bson:"first_name" json:"first_name"`
	Last_name   string   `bson:"last_name" json:"last_name"`
	Email       *string  `bson:"email" json:"email"`
	Phone       *string  `bson:"phone" json:"phone"`
	Description *string  `bson:"description" json:"description"`
	Image       *string  `bson:"image" json:"image"`
	Rating      *float64 `bson:"rating" json:"rating"`
}
