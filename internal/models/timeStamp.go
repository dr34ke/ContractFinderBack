package models

import (
	"time"
)

type TimeStamp struct {
	Created_at time.Time
	Updated_at time.Time
}

func (t TimeStamp) Created() {
	t.Created_at = time.Now()
	t.Updated_at = time.Now()
}

func (t TimeStamp) Updated() {
	t.Updated_at = time.Now()
}
