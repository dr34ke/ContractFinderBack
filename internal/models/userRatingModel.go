package models

type UserRating struct {
	Id              string `bson:"_id" json:"id"`
	IdEvaluatorUser string `bson:"idEvaluatorUser" json:"ideEvaluatorUser"`
	IdEvaluatedUser string `bson:"idEvaluatedUser" json:"ideEvaluatedUser"`
	AsWorker        bool   `bson:"asWorker" json:"asWorker"`
	Rating          string `bson:"rating" json:"rating"`
	Description     string `bson:"description" json:"description"`
	TimeStamp       TimeStamp
}
