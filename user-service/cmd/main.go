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
	"github.com/winkedin/user-service/models"
	"github.com/winkedin/user-service/services"
)

var (
	configFilePath = flag.String("config", "", "absolute path to the config file")
)

func main() {
	flag.Parse()
	v, err := services.GetConfig(*configFilePath)
	if err != nil {
		panic(err)
	}
	db, err := services.GetDBConnection(v, &models.User{})
	if err != nil {
		panic(err)
	}
	rdb, err := services.GetRedisConnection(context.Background(), *v)
	if err != nil {
		panic(err)
	}
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
		err := emailVerificationSVC.VerifyOTP(context.Background(), userVerifyRequest.Email, userVerifyRequest.OTP)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		rdb.HGetAll(context.Background(), fmt.Sprintf("user:%s", userVerifyRequest.Email)).Scan(&user)
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
		token, err := loginSvc.Login(context.Background(), userLoginRequest.Email, userLoginRequest.LinkedInJWT)
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
		err := loginSvc.Logout(context.Background(), userLogoutRequest.UserID)
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
