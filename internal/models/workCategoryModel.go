package models

type WorkCategory struct {
	Id          string `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	OffersCount int    `bson:"offersCount,omitempty" json:"offersCount,omitempty"`
}
