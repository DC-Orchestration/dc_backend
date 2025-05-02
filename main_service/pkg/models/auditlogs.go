package models

import "time"

type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DocID     uint      `json:"doc_id"`
	Action    string    `json:"action"`
	UserID    uint      `json:"user_id"`
	Hash      string    `json:"hash"`
	PrevHash  string    `json:"prev_hash"`
	Summary   string    `json:"summary"`
	Timestamp time.Time `json:"timestamp"`
}
