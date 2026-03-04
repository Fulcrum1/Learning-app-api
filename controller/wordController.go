package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"learning-app-mobile-bna/Backend/initializers"
	"learning-app-mobile-bna/Backend/model"
)

type Word struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Translation   string `json:"translation"`
	Pronunciation string `json:"pronunciation"`
}

type CreateWordRequest struct {
	Word          string `json:"word" binding:"required"`
	Translation   string `json:"translation" binding:"required"`
	Pronunciation string `json:"pronunciation"`
	Language      string `json:"language" binding:"required"`
}

type CreateWordsRequest struct {
	Words    []CreateWordRequest `json:"words" binding:"required"`
	Language string              `json:"language" binding:"required"`
}

type UpdateWordRequest struct {
	Word          string `json:"word" binding:"required"`
	Translation   string `json:"translation" binding:"required"`
	Pronunciation string `json:"pronunciation"`
}

// Get all words from user in specific language
func GetWords(c *gin.Context) {
	var words []model.Word
	language := GetLanguageId(c.Query("language"))

	if language == 0 {
		initializers.DB.Find(&words)
	} else {
		initializers.DB.Where("language_id = ?", language).Find(&words)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": words,
	})
}

func GetWord(c *gin.Context) {
	wordId := c.Param("id")

	var word model.Word
	initializers.DB.First(&word, wordId)

	c.JSON(http.StatusOK, gin.H{
		"data": word,
	})
}

func createSingleWord(word, translation, pronunciation, language string) (int64, string, int) {
	if word == "" || translation == "" {
		return 0, "Le mot et sa traduction sont obligatoires", 400
	}

	languageId := GetLanguageId(language)
	if languageId == 0 {
		return 0, "Langue non supportée", 400
	}

	exists, err := checkWordExists(word, languageId)
	if err != nil {
		return 0, "Erreur lors de la vérification du mot", 500
	}
	if exists {
		return 0, "Ce mot existe déjà dans votre vocabulaire", 409
	}

	wordID, err := AddWordToDatabase(word, translation, pronunciation, languageId)
	if err != nil {
		return 0, "Erreur lors de l'ajout du mot", 500
	}

	return int64(wordID), "", 201
}

func CreateWord(c *gin.Context) {
	var req CreateWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	wordID, errMsg, status := createSingleWord(req.Word, req.Translation, req.Pronunciation, req.Language)
	if status != 201 {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(201, gin.H{
		"message":       "Mot créé avec succès",
		"word_id":       wordID,
		"word":          req.Word,
		"translation":   req.Translation,
		"pronunciation": req.Pronunciation,
	})
}

func CreateWords(c *gin.Context) {
	var req CreateWordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var results []gin.H
	for _, word := range req.Words {
		wordID, errMsg, status := createSingleWord(word.Word, word.Translation, word.Pronunciation, req.Language)
		if status != 201 {
			results = append(results, gin.H{
				"word":  word.Word,
				"error": errMsg,
			})
		} else {
			results = append(results, gin.H{
				"word_id":       wordID,
				"word":          word.Word,
				"translation":   word.Translation,
				"pronunciation": word.Pronunciation,
			})
		}
	}

	c.JSON(201, gin.H{
		"message": "Mots créés avec succès",
		"results": results,
	})
}

// Update a word for user in specific language
func UpdateWord(c *gin.Context) {
	wordId := c.Param("id")
	var word UpdateWordRequest

	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	initializers.DB.Model(&model.Word{}).Where("id = ?", wordId).Updates(word)
	c.JSON(200, gin.H{
		"message": "Word updated successfully",
	})
}

// Delete a word for user in specific language
func DeleteWord(c *gin.Context) {
	wordId := c.Param("id")
	var listsWord []model.ListWord

	// Get all lists that contain this word
	initializers.DB.Where("word_id = ?", wordId).Find(&listsWord)
	// Delete the word from all lists
	initializers.DB.Where("word_id = ?", wordId).Delete(&model.ListWord{})
	// Delete the word
	initializers.DB.Delete(&model.Word{}, wordId)

	for _, listWord := range listsWord {
		var listWordModel model.ListWord
		initializers.DB.Where("list_id = ?", listWord.ListID).Find(&listWordModel)

		if listWordModel.ID == 0 {
			initializers.DB.Delete(&model.List{}, listWord.ListID)
		}
	}

	c.JSON(200, gin.H{
		"data": "success",
	})
}

func GetCountWords(c *gin.Context) int64 {
	var count int64

	languageId := GetLanguageId(c.Query("language"))
	userId := c.Query("user")

	if languageId == 0 {
		c.JSON(400, gin.H{
			"error": "Language not found",
		})
		return 0
	}

	initializers.DB.Model(&model.Word{}).Where("language_id = ? AND user_id = ?", languageId, userId).Count(&count)

	return count
}

// Retourne []model.Word
func FetchKnownWords(languageId uint, userId string) []model.Word {
	words := []model.Word{}
	initializers.DB.Where("language_id = ? AND user_id = ? AND score >= 80", languageId, userId).Find(&words)
	return words
}

// Retourne []model.Word
func FetchLearningWords(languageId uint, userId string) []model.Word {
	words := []model.Word{}
	initializers.DB.Where("language_id = ? AND user_id = ? AND score > 20 AND score < 80", languageId, userId).Find(&words)
	return words
}

// Retourne []model.Word
func FetchUnknownWords(languageId uint, userId string) []model.Word {
	words := []model.Word{}
	initializers.DB.Where("language_id = ? AND user_id = ? AND (score <= 20 OR score IS NULL)", languageId, userId).Find(&words)
	return words
}

// ------------------------------------------------------------------------------------------
// ------------------------------------ Utils Functions ------------------------------------
// ------------------------------------------------------------------------------------------

func checkWordExists(word string, languageId uint) (bool, error) {
	var wordExists model.Word
	initializers.DB.Where("word = ? AND language_id = ?", word, languageId).First(&wordExists)

	return wordExists.ID != 0, nil
}

// func getLanguageId(langague string) uint {
// 	var language model.Language
// 	initializers.DB.Where("LOWER(code) = LOWER(?)", langague).First(&language)
// 	return language.ID
// }

func AddWordToDatabase(word string, translation string, pronunciation string, language uint) (uint, error) {
	var wordAdd model.Word

	wordAdd.Word = word
	wordAdd.Translation = translation
	wordAdd.Pronunciation = pronunciation
	wordAdd.LanguageID = language

	initializers.DB.Create(&wordAdd)

	return wordAdd.ID, nil
}

// func addWordToList(c *gin.Context) {

// }
// func initializeWordProgress(c *gin.Context) {

// }
