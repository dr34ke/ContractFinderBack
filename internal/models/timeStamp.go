package models

import (
	"time"
)

type TimeStamp struct {
	Created_at time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	Updated_at time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	Last_login time.Time `json:"lastLogin,omitempty" bson:"lastLogin,omitempty"`
}

func (t *TimeStamp) Created() {
	t.Created_at = time.Now()
}

func (t *TimeStamp) Updated() {
	t.Updated_at = time.Now()
}

func (t *TimeStamp) Login() {
	t.Last_login = time.Now()
}
