package model

import "time"

type Word struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Word          string `json:"word"`
	Translation   string `json:"translation"`
	Pronunciation string `json:"pronunciation"`
	Score         uint   `json:"score" max:"100" min:"0"`
	Review        bool   `json:"review"`

	UserID     uint `json:"user_id"`
	LanguageID uint `json:"language_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
