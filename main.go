package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"learning-app-mobile-bna/Backend/controller"
	"learning-app-mobile-bna/Backend/initializers"
	"learning-app-mobile-bna/Backend/middleware"
	"learning-app-mobile-bna/Backend/model"

	"github.com/gin-contrib/cors"
)

func init() {
	initializers.LoadEnvs()
	initializers.ConnectDB()
}

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()

	router.Use(gin.Logger())
	// router.Use(cors.Default())
	allowedOrigins := []string{"http://localhost:3000"}
	if prodOrigin := os.Getenv("ALLOWED_ORIGIN"); prodOrigin != "" {
		allowedOrigins = append(allowedOrigins, prodOrigin)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Define a simple GET endpoint
	router.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Auth routes
	router.POST("/auth/login", controller.Login)
	router.POST("/auth/register", controller.CreateUser)

	protected := router.Group("/api")
	protected.Use(middleware.AuthenticateMiddleware)

	protected.GET("/users", controller.GetUsers)

	// Languages routes
	protected.POST("/languages", controller.CreateLanguage)

	protected.GET("/dashboard", func(c *gin.Context) {
		var totalWords = controller.GetCountWords(c)

		languageId := controller.GetLanguageId(c.Query("language"))
		userId := c.Query("user")

		var completed []model.Word
		var mastered []model.Word
		var learning []model.Word
		var toReview []model.Word

		initializers.DB.Where("user_id = ? AND language_id = ?", userId, languageId).Find(&completed)
		mastered = controller.FetchKnownWords(languageId, userId)
		learning = controller.FetchLearningWords(languageId, userId)
		toReview = controller.FetchUnknownWords(languageId, userId)

		c.JSON(http.StatusOK, gin.H{
			"totalWords":      totalWords,
			"completed":       completed,
			"mastered":        mastered,
			"learning":        learning,
			"toReview":        toReview,
			"lastListLearned": controller.GetLastListLearned(languageId, userId),
		})
	})

	// Words routes
	protected.GET("/words", controller.GetWords)
	protected.GET("/word/:id", controller.GetWord)
	protected.POST("/word", controller.CreateWord)
	protected.POST("/words", controller.CreateWords)
	protected.PUT("/word/:id", controller.UpdateWord)
	protected.DELETE("/word/:id", controller.DeleteWord)
	protected.GET("/words/count", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{"count": controller.GetCountWords(c)},
		)
	})

	// Lists routes
	protected.GET("/lists", controller.GetLists)
	protected.POST("/list", controller.CreateList)
	protected.GET("/list/:id", controller.GetList)
	protected.PUT("/list/:id", controller.UpdateList)
	protected.DELETE("/list/:id", controller.DeleteList)
	protected.PUT("/review-end/:id", controller.ReviewEnd)

	protected.GET("/card/words-card/:id", controller.GetWordsCard)
	protected.GET("/card/param-card/:id", controller.GetCardParam)
	protected.POST("/card/progress-card/:id", controller.ProgressCard)
	protected.POST("/card/rollback-progress-card", controller.RollbackProgressCard)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
