package model

type AuthInput struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"required"`
}
