package model

import "time"

type List struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`

	DateReview time.Time `json:"date_review"`
	CountWords uint      `json:"count_words"`

	UserID     uint `json:"user_id"`
	LanguageID uint `json:"language_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ListWord []ListWord `json:"list_word"`
}
