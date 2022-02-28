package models

import (
	"time"
)

type Post struct {
	ID			uint `gorm:"primaryKey" json:"id"`
	UserID		uint `json:"-"`
	Email	string `json:"email"`
	Content		string `json:"content"`
	Create_At	time.Time `json:"create_at"`	
}
