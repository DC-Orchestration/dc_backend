package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"password"` // Should be hashed
	Signature string    `json:"signature"`
	CreatedAt time.Time `json:"created_at"`
	Groups    []Group   `gorm:"many2many:group_users;" json:"groups"`
}
