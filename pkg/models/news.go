package models

// News model.
type News struct {
	CreatedUpdated

	ID     int64  `json:"id" db:"id"`
	Author string `json:"author" db:"author" valid:"required"`
	Body   string `json:"body" db:"body" valid:"required"`
}
