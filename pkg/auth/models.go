package auth

import "time"

type User struct {
	ID          string
	DisplayName string
	Email       string
	CreatedAt   time.Time
}
