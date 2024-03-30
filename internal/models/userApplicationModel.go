package models

type UserApplication struct {
	Id          string  `bson:"_id" json:"id"`
	UserId      string  `bson:"userId" json:"userId"`
	OfferId      string  `bson:"offerId" json:"offerId"`
	SalaryOffer float32 `bson:"salaryOffer" json:"salaryOffer"`
	Description string  `bson:"description" json:"description"`
}
