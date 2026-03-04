package model

type ListWord struct {
	ID     uint `json:"id"`
	ListID uint `json:"list_id"`
	WordID uint `json:"word_id"`
	Review bool `json:"review"`

	Word Word `gorm:"foreignKey:WordID;references:ID" json:"word"`
	// List List `gorm:"foreignKey:ListID;references:ID" json:"list"`
}
