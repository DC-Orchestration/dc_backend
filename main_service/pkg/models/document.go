package models

import "time"

type Document struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	GroupID    uint       `gorm:"not null" json:"group_id"`
	UploaderID uint       `gorm:"not null" json:"uploader_id"`
	DocName    string     `gorm:"size:255;not null" json:"doc_name"`
	Content    string     `gorm:"type:text" json:"content"`
	Version    int        `gorm:"default:1" json:"version"`
	Hash       string     `gorm:"not null;unique" json:"hash"`   
	PrevHash   string     `gorm:"not null" json:"prev_hash"`  
	Path string `gorm:"not null" json:"path"`
   
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	AuditLogs []AuditLog `gorm:"foreignKey:DocID;constraint:OnDelete:CASCADE" json:"audit_logs"`
}
