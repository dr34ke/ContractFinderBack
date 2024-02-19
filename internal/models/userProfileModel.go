package models

type UserProfile struct {
	Description string `bson:"description" json:"description"`
	Image       string `bson:"image" json:"image"`
	TimeStamp   TimeStamp
}
