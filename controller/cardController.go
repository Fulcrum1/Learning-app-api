package controller

import (
	"learning-app-mobile-bna/Backend/initializers"
	"learning-app-mobile-bna/Backend/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetWordsCard(c *gin.Context) {
	var listWords []model.ListWord
	var list model.List

	checkReviewList(c.Param("id"))

	initializers.DB.Where("id = ?", c.Param("id")).First(&list)

	initializers.DB.Where("list_id = ? AND (review = false OR review = null)", c.Param("id")).Preload("Word").Find(&listWords)

	count := len(listWords)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"words": listWords,
		"count": count,
	})
}

func GetCardParam(c *gin.Context) {
	var params model.Params
	userId := c.Param("id")

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	initializers.DB.Where("user_id = ?", userId).First(&params)

	if params.ID == 0 {
		params = model.Params{
			UserID:             uint(userIdInt),
			Random:             false,
			TranslationOnVerso: false,
		}
		initializers.DB.Create(&params)
	}

	c.JSON(200, gin.H{"data": params})
}

func ProgressCard(c *gin.Context) {
	listId := c.Param("id")

	var body struct {
		WordId  any  `json:"wordId"`
		IsKnown bool `json:"isKnown"`
		UserId  any  `json:"userId"`
	}
	c.BindJSON(&body)

	listWord := model.ListWord{}
	word := model.Word{}

	initializers.DB.Where("list_id = ? AND word_id = ?", listId, body.WordId).First(&listWord)
	initializers.DB.Where("user_id = ? AND id = ?", body.UserId, body.WordId).First(&word)

	score := 20
	if body.IsKnown {
		remaining := 100 - int(word.Score)
		gain := (score * remaining) / 100
		if gain < 1 {
			gain = 1
		}
		word.Score = uint(min(int(word.Score)+gain, 100))
		initializers.DB.Save(&word)
		listWord.Review = true
	} else {
		if int(word.Score)-score < 0 {
			word.Score = 0
		} else {
			word.Score -= uint(score)
		}
		initializers.DB.Save(&word)
		listWord.Review = false
	}

	initializers.DB.Save(&listWord)
	checkReviewList(listId)

	c.JSON(200, gin.H{"data": word})
}

func RollbackProgressCard(c *gin.Context) {

}

func checkReviewList(listId string) {
	var listWords []model.ListWord
	initializers.DB.Where("list_id = ?", listId).Find(&listWords)

	var reviewFalse bool = false
	for _, el := range listWords {
		if el.Review == false {
			reviewFalse = true
		}
	}

	if !reviewFalse {
		initializers.DB.Model(&model.ListWord{}).Where("list_id = ?", listId).Update("review", false)
	}
}
