package models

type UserApplication struct {
	Id          string  `bson:"_id" json:"id"`
	UserId      string  `bson:"userId" json:"userId"`
	WorkId      string  `bson:"workId" json:"workId"`
	SalaryOffer float32 `bson:"salaryOffer" json:"salaryOffer"`
	Description string  `bson:"description" json:"description"`
}
