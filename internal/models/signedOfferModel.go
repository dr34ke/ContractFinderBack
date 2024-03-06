package models

type SignedOfferModel struct {
	Id        string    `bson:"_id" json:"id"`
	UserId    string    `bson:"userId" json:"userId"`
	WorkId    string    `bson:"workId" json:"workId"`
	TimeStamp TimeStamp `bson:"timeStamp" json:"timeStamp"`
}
