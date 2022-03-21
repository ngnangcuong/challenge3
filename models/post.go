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

type PostRepo interface {
	Select() ([]Post, error)
	Delete(id string) (error)
	Update(id string, content string)	(error)
	Create(post Post) (error)
	Find(id string) (Post, error)
}