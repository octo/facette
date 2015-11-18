package backend

import (
	"time"
)

type Item struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Created     time.Time `db:"created" json:"created"`
	Modified    time.Time `db:"modified" json:"modified"`
}
