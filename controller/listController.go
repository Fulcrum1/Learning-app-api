package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"learning-app-mobile-bna/Backend/initializers"
	"learning-app-mobile-bna/Backend/model"
	// "learning-app-mobile-bna/Backend/controller/languageController"
)

type CreateListRequest struct {
	Name     string `json:"name"`
	UserId   uint   `json:"userId"`
	Language string `json:"language"`
	Words    []uint `json:"words"`
}

type UpdateListRequest struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	UserId   uint   `json:"user_id"`
	Language string `json:"language"`
	Words    []uint `json:"words"`
}

type AddWordToListRequest struct {
	ListId uint `json:"list_id"`
	WordId uint `json:"word_id"`
	Review bool `json:"review"`
}

func GetLists(c *gin.Context) {
	var lists []model.List
	language := GetLanguageId(c.Query("language"))
	userId := c.Query("user_id")

	if language == 0 {
		initializers.DB.Where("user_id = ?", userId).Find(&lists)
	} else {
		initializers.DB.Where("user_id = ? AND language_id = ?", userId, language).Find(&lists)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": lists,
	})
}

func GetList(c *gin.Context) {
	var listWords []model.ListWord
	var list model.List

	initializers.DB.Where("id = ?", c.Param("id")).First(&list)

	initializers.DB.Where("list_id = ?", c.Param("id")).Preload("Word").Find(&listWords)

	count := len(listWords)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"words": listWords,
		"count": count,
	})
}

func CreateList(c *gin.Context) {
	var req CreateListRequest
	fmt.Println(req)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Name == "" || req.UserId == 0 || req.Language == "" {
		c.JSON(400, gin.H{
			"error": "Le nom, l'utilisateur et la langue sont obligatoires",
		})
		return
	}

	languageId := GetLanguageId(req.Language)
	if languageId == 0 {
		c.JSON(400, gin.H{
			"error": "Langue non supportée",
		})
		return
	}

	countWords := uint(len(req.Words))

	listAdd := AddListToDatabase(req, languageId, countWords)

	for _, wordId := range req.Words {
		AddWordToList(
			AddWordToListRequest{
				ListId: listAdd.ID,
				WordId: wordId,
			}, false)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": listAdd,
	})
}

func UpdateList(c *gin.Context) {
	var req UpdateListRequest

	list := model.List{}
	initializers.DB.Where("id = ?", req.Id).Find(&list)

	list.Name = req.Name

	initializers.DB.Save(&list)

	var listWord model.ListWord
	initializers.DB.Where("list_id = ?", list.ID).Delete(&listWord)

	for _, wordId := range req.Words {
		AddWordToList(
			AddWordToListRequest{
				ListId: list.ID,
				WordId: wordId,
			}, false)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

func DeleteList(c *gin.Context) {
	var list model.List
	initializers.DB.Where("id = ?", c.Param("id")).Delete(&list)

	var listWord model.ListWord
	initializers.DB.Where("list_id = ?", list.ID).Delete(&listWord)

	c.JSON(http.StatusOK, gin.H{
		"message": "List deleted successfully",
	})
}

func ReviewEnd(c *gin.Context) {
	var list model.List
	var listWords []model.ListWord

	initializers.DB.Where("id = ?", c.Param("id")).Find(&list)
	if err := initializers.DB.Preload("Word").Where("list_id = ?", c.Param("id")).Find(&listWords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch list words",
		})
		return
	}

	// for _, listWord := range listWords {
	// 	UpdateWordProgressToDatabase(listWord.WordID, listWord.Word.UserID, false, listWord.Word.Score)
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Review ended successfully",
		"data":    listWords,
	})
}

func GetLastListLearned(languageId uint, userId string) model.List {
	var list model.List
	initializers.DB.Where("user_id = ?", userId).Preload("ListWord").Order("updated_at desc").First(&list)
	return list
}

// ---------------------- Utils ----------------------
func AddListToDatabase(req CreateListRequest, languageId uint, countWords uint) model.List {
	var listAdd model.List

	listAdd.Name = req.Name
	listAdd.UserID = req.UserId
	listAdd.LanguageID = languageId
	listAdd.CountWords = countWords

	initializers.DB.Create(&listAdd)

	return listAdd
}

func AddWordToList(req AddWordToListRequest, review bool) {
	listWord := model.ListWord{}
	listWord.ListID = req.ListId
	listWord.WordID = req.WordId
	listWord.Review = review

	initializers.DB.Create(&listWord)
}

func DeleteWordFromList(req AddWordToListRequest) {
	var listWord model.ListWord
	initializers.DB.Where("list_id = ? AND word_id = ?", req.ListId, req.WordId).Delete(&listWord)
}
