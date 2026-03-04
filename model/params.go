package model

type Params struct {
	ID                 uint `json:"id"`
	Random             bool `json:"random"`
	TranslationOnVerso bool `json:"translation_on_verso"`

	UserID uint `json:"user_id"`
}
