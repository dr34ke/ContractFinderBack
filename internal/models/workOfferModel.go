package models

type WorkOffer struct {
	Id                string            `bson:"_id" json:"id"`
	UserId            string            `bson:"userId" json:"userId" validate:"required"`
	CategoryId        string            `bson:"categoryId" json:"categoryId" validate:"required"`
	Title             string            `bson:"title" json:"title" validate:"required"`
	Description       string            `bson:"description" json:"description" validate:"required"`
	SugestedSalary    float32           `bson:"sugestedSalary" json:"sugestedSalary" validate:"required"`
	IsSalaryPerHour   bool              `bson:"isSalaryPerHour" json:"isSalaryPerHour"`
	IsRepetetive      bool              `bson:"isRepetetive" json:"isRepetetive"`
	IsFromWorker      bool              `bson:"isFromWorker" json:"isFromWorker"`
	OnSite            bool              `bson:"onSite" json:"onSite"`
	Coordinates       []float64         `bson:"coordinates" json:"coordinates"`
	DistanceInKm      float64           `bson:"distanceInKm,omitempty" json:"distanceInKm"`
	UsersApplications []UserApplication `bson:"usersApplications,omitempty" json:"usersApplications"`
	TimeStamp         TimeStamp         `bson:"timeStamp,omitempty" json:"timeStamp"`
}
