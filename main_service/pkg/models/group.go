package models

import "time"

type Group struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Signature string    `json:"signature"`
	Members   []User    `gorm:"many2many:group_users;" json:"members"`
	Documents []Document `gorm:"foreignKey:GroupID" json:"documents"`
}
