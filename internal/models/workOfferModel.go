package models

type WorkOffer struct {
	Id                string            `bson:"_id" json:"id"`
	UserId            string            `bson:"userId" json:"userId"`
	CategoryId        string            `bson:"categoryId" json:"categoryId"`
	Title             string            `bson:"title" json:"title"`
	Description       string            `bson:"description" json:"description"`
	SugestedSalary    float32           `bson:"sugestedSalary" json:"sugestedSalary"`
	IsSalaryPerHour   bool              `bson:"isSalaryPerHour" json:"isSalaryPerHour"`
	IsRepetetive      bool              `bson:"isRepetetive" json:"isRepetetive"`
	IsFromWorker      bool              `bson:"isFromWorker" json:"isFromWorker"`
	OnSite            bool              `bson:"onSite" json:"onSite"`
	Coordinates       []float64         `bson:"coordinates" json:"coordinates"`
	DistanceInKm      float64           `bson:"distanceInKm" json:"distanceInKm"`
	UsersApplications []UserApplication `bson:"usersApplications" json:"usersApplications"`
	TimeStamp         TimeStamp         `bson:"timeStamp" json:"timeStamp"`
}
