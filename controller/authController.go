package controller

import (
	"net/http"

	"learning-app-mobile-bna/Backend/initializers"
	"learning-app-mobile-bna/Backend/model"

	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {

	var authInput model.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound model.User
	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already used"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		Email:    authInput.Email,
		Name:     authInput.Name,
		Password: string(passwordHash),
	}

	initializers.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func Login(c *gin.Context) {

	var authInput model.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound model.User
	// Utilisez First() au lieu de Find()
	result := initializers.DB.Where("LOWER(email) = LOWER(?)", authInput.Email).First(&userFound)

	// Vérifiez l'erreur retournée par First()
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
		return // N'oubliez pas le return ici !
	}

	c.JSON(200, gin.H{
		"token": token,
		"user":  userFound,
	})
}

func GetUserProfile(c *gin.Context) {

	user, _ := c.Get("currentUser")

	c.JSON(200, gin.H{
		"user": user,
	})
}
