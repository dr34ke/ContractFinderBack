package models

type UserPreference struct {
	IsPhonePublic bool      `bson:"isPhonePublic" json:"isPhonePublic"`
	IsEmailPublic bool      `bson:"isEmailPublic" json:"isEmailPublic"`
	UserType      string    `bson:"userType" json:"userType"`
	WorkDistance  float32   `bson:"workDistance" json:"workDistance,string"`
	TimeStamp     TimeStamp `bson:"timeStamp" json:"timeStamp"`
}
