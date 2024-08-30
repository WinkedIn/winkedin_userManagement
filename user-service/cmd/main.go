package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"github.com/winkedin/user-service/models"
	"github.com/winkedin/user-service/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	configFilePath = flag.String("config", "", "absolute path to the config file")
)

var (
	db  *gorm.DB
	rdb *redis.Client
	ctx = context.Background()
	v   = viper.New()
)

func initConfig() {
	v.SetConfigFile(*configFilePath)
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func initDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		v.GetString("database.host"),
		v.GetString("database.user"),
		v.GetString("database.password"),
		v.GetString("database.dbname"),
		v.GetInt("database.port"),
		v.GetString("database.sslmode"),
	)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	db.AutoMigrate(&models.User{})
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", v.GetString("redis.host"), v.GetInt("redis.port")),
		Password: v.GetString("redis.password"),
		DB:       0,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}
}

func main() {
	flag.Parse()
	initConfig()
	initDB()
	initRedis()
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Println(fmt.Sprintf("%v", err))
		}
	}()
	emailVerificationSVC := services.NewEmailVerificationService(rdb)
	signupSVC := services.NewSignupService(db, rdb, emailVerificationSVC)
	loginSvc := services.NewLoginService(db, rdb)
	router := gin.Default()

	router.POST("/signup", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if err := signupSVC.Signup(c.Request.Context(), &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Signup initiated. Please verify your email."})
	})

	router.POST("/verify", func(c *gin.Context) {
		var userVerifyRequest models.VerifyRequest
		if err := c.ShouldBindJSON(&userVerifyRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		err := emailVerificationSVC.VerifyOTP(ctx, userVerifyRequest.Email, userVerifyRequest.OTP)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		rdb.HGetAll(ctx, fmt.Sprintf("user:%s", userVerifyRequest.Email)).Scan(&user)
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user data"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Account verified and created successfully!"})
	})

	router.POST("/login", func(c *gin.Context) {
		var userLoginRequest models.UserLoginRequest
		if err := c.ShouldBindJSON(&userLoginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		token, err := loginSvc.Login(ctx, userLoginRequest.Email, userLoginRequest.LinkedInJWT)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
	})

	router.POST("/logout", func(c *gin.Context) {
		var userLogoutRequest models.UseLogoutRequest
		if err := c.ShouldBindJSON(&userLogoutRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		err := loginSvc.Logout(ctx, userLogoutRequest.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	})
	go func() {
		if err := router.Run(fmt.Sprintf(":%d", v.GetInt("app.port"))); err != nil {
			panic(err)
		}
	}()
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGTERM, syscall.SIGINT)
	<-quitSignal
}
