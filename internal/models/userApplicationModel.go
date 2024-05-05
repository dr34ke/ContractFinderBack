package models

type UserApplication struct {
	Id          string  `bson:"_id" json:"id"`
	UserId      string  `bson:"userId" json:"userId" validate:"required"`
	OfferId     string  `bson:"offerId" json:"offerId" validate:"required"`
	SalaryOffer float32 `bson:"salaryOffer" json:"salaryOffer" validate:"required"`
	Description string  `bson:"Description" json:"description" validate:"required"`
}
