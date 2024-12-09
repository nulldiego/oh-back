package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/nulldiego/oh-back/config"
	"github.com/nulldiego/oh-back/internal/database"
	"github.com/nulldiego/oh-back/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var reqBody struct {
		Username string `json:"username" binding:"required,min=5"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var count int64
	if err := database.DB.Table("users").Where("username = ?", reqBody.Username).Count(&count).Error; err != nil || count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already taken"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	user := model.User{
		Username: reqBody.Username,
		Password: string(hashPassword),
	}

	if err = database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Login(c *gin.Context) {
	var reqBody struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := database.DB.First(&user, "username = ?", reqBody.Username).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func Logout(c *gin.Context) {
	// Clear the cookie
	c.SetCookie("Authorization", "", 0, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func generateToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Conf.JwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
