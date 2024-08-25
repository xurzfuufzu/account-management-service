package models

import (
	"time"
)

type User struct {
	ID        string
	Username  string
	Password  string
	CreatedAt time.Time
}
