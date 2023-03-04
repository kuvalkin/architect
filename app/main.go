package main

import (
	"net/http"
	"strconv"
	"time"

	"errors"

	"os"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	"gorm.io/driver/postgres"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

var METRICS_REQUEST_LATENCY = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "app_request_latency_seconds",
		Help: "Application Request Latency",
	},
	[]string{"method", "endpoint"},
)

var METRICS_REQUEST_COUNT = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "app_request_count",
		Help: "Application Request Count",
	},
	[]string{"method", "endpoint", "http_status"},
)

func main() {
	r := gin.Default()
	db := initDb()

	// define middleware before all routes, otherwise it will hurt you :)
	r.Use(func(c *gin.Context) {

		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()

		c.Next()

		ellapsed := float64(time.Since(start)) / float64(time.Second)

		METRICS_REQUEST_LATENCY.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(ellapsed)
		METRICS_REQUEST_COUNT.WithLabelValues(c.Request.Method, c.Request.URL.Path, strconv.Itoa(c.Writer.Status()))
	})

	promHandler := promhttp.Handler()
	r.GET("/metrics", func(c *gin.Context) {
		promHandler.ServeHTTP(c.Writer, c.Request)
	})

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

	r.PUT("/user/:userId", func(c *gin.Context) {
		var user UserModel
		result := db.First(&user, c.Param("userId"))

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"code": 3, "message": "Not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 4, "message": result.Error.Error()})
			}

			return
		}

		var payload UserInput
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 5, "message": err.Error()})
			return
		}

		result = db.Model(&user).Updates(&UserModel{
			Username:  payload.Username,
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Email:     payload.Email,
			Phone:     payload.Phone,
		})
		if result.Error == nil {
			c.JSON(http.StatusOK, user)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"code": 6, "message": result.Error.Error()})
		}
	})

	r.DELETE("/user/:userId", func(c *gin.Context) {
		result := db.Delete(&UserModel{}, c.Param("userId"))

		if result.Error == nil {
			c.Status(http.StatusNoContent)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 7, "message": result.Error.Error()})
		}
	})

	r.Run(":8000")
}

func initDb() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&UserModel{})

	return db
}
