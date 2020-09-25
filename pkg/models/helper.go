package models

import (
	"time"
)

// CreatedUpdated stores created and updated time of the model.
type CreatedUpdated struct {
	Created *time.Time `json:"created,omitempty" db:"created" example:"2020-04-21T00:00:00Z"`
}
