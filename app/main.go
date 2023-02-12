package main

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	"gorm.io/driver/sqlite"
)

type UserModel struct {
	gorm.Model
	Username  string
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

type UserInput struct {
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"` //todo validation
	Phone     string `json:"phone" binding:"required"` //todo validation
}

func main() {
	r := gin.Default()
	db := initDb()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	r.POST("/user", func(c *gin.Context) {
		var payload UserInput
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
			return
		}

		user := &UserModel{
			Username:  payload.Username,
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Email:     payload.Email,
			Phone:     payload.Phone,
		}

		result := db.Create(user)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 2, "message": result.Error.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	})

	r.GET("/user/:userId", func(c *gin.Context) {
		var user UserModel
		result := db.First(&user, c.Param("userId"))
		if result.Error == nil {
			c.JSON(http.StatusOK, user)
			return
		}

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 3, "message": "Not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 4, "message": result.Error.Error()})
		}
	})

	r.Run(":8000")
}

func initDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&UserModel{})

	return db
}
