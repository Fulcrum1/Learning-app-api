package controller

import (
	"learning-app-mobile-bna/Backend/initializers"
	"learning-app-mobile-bna/Backend/model"

	"github.com/gin-gonic/gin"

	"net/http"
)

func CreateLanguage(c *gin.Context) {
	name := c.Query("name")
	code := c.Query("code")

	initializers.DB.Create(&model.Language{
		Name: name,
		Code: code,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Language created successfully",
	})
}

func GetLanguageId(langague string) uint {
	var language model.Language
	initializers.DB.Where("LOWER(code) = LOWER(?)", langague).First(&language)
	return language.ID
}
