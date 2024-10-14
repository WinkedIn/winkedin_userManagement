package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/winkedin/user-service/store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/winkedin/user-service/constants"
	wlog "github.com/winkedin/user-service/logger"
	"github.com/winkedin/user-service/models"
	"github.com/winkedin/user-service/services"
)

// extractRequestIdAndBuildContext extracts the request ID from the gin context and adds it to the request context
func extractRequestIdAndBuildContext(c *gin.Context) context.Context {
	requestID := requestid.Get(c)
	clientIp := c.ClientIP()
	// build request path
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}
	requestPath := path
	requestMethod := c.Request.Method
	ctx := context.WithValue(c.Request.Context(), constants.Key(constants.RequestIdKey), requestID)
	ctx = context.WithValue(ctx, constants.Key(constants.UserIP), clientIp)
	ctx = context.WithValue(ctx, constants.Key(constants.RequestPath), requestPath)
	ctx = context.WithValue(ctx, constants.Key(constants.RequestMethod), requestMethod)
	return ctx
}

func main() {
	flag.Parse()
	v := services.GetConfig(*services.ConfigFilePath)
	db, err := store.GetDBConnection(v, &models.User{})
	if err != nil {
		panic(err)
	}
	rdb, err := services.GetRedisConnection(context.Background(), *v)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("%v", err)
		}
	}()
	wlog.InitLogger()
	emailVerificationSVC := services.NewEmailVerificationService(rdb)
	signupSVC := services.NewSignupService(db, rdb, emailVerificationSVC)
	signInWithLinkedInSvc := services.NewSignInWithLinkedInService(db, rdb)
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(logger.SetLogger())
	router.Use(requestid.New())

	router.POST("/signup", func(c *gin.Context) {
		ctx := extractRequestIdAndBuildContext(c)
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if err := signupSVC.Signup(ctx, &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Signup initiated. Please verify your email."})
	})

	router.POST("/verify", func(c *gin.Context) {
		ctx := extractRequestIdAndBuildContext(c)
		var userVerifyRequest models.VerifyRequest
		if err = c.ShouldBindJSON(&userVerifyRequest); err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		err = emailVerificationSVC.VerifyOTP(ctx, userVerifyRequest.Email, userVerifyRequest.OTP)
		if err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		err = rdb.HGetAll(ctx, fmt.Sprintf("user:%s", userVerifyRequest.Email)).Scan(&user)
		if err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data from Redis"})
			return
		}
		if err = db.Create(&user).Error; err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user data"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Account verified and created successfully!"})
	})

	// group callback routes
	callback := router.Group("/callback")

	// group OAuth callbacks
	oAuthRoutes := callback.Group("/oauth")
	// LinkedIn OAuth callback
	oAuthRoutes.POST("/linkedin", func(c *gin.Context) {
		ctx := extractRequestIdAndBuildContext(c)
		// fetch LinkedIn auth code from query params
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auth code"})
			return
		}
		// get LinkedIn profile and login
		token, err := signInWithLinkedInSvc.GetLinkedInProfileAndLogin(ctx, code)
		if err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
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
