package models

type UserRating struct {
	Id              string    `bson:"_id" json:"id"`
	IdEvaluatorUser string    `bson:"idEvaluatorUser" json:"idEvaluatorUser"`
	IdEvaluatedUser string    `bson:"idEvaluatedUser" json:"idEvaluatedUser"`
	AsWorker        bool      `bson:"asWorker" json:"asWorker"`
	Rating          int32     `bson:"rating" json:"rating"`
	Description     string    `bson:"description" json:"description"`
	TimeStamp       TimeStamp `bson:"timeStamp" json:"timeStamp"`
}
