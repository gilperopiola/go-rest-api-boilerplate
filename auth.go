package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	Email          string
	Password       string
	RepeatPassword string
	Code           string
}

func validateToken(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")

		if len(tokenString) < 40 {
			c.JSON(http.StatusUnauthorized, "authentication error")
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.SECRET), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, "authentication error")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			c.Set("ID", claims.Id)
			c.Set("Email", claims.Audience)
		} else {
			c.JSON(http.StatusUnauthorized, "authentication error")
			c.Abort()
		}
	}
}

func generateToken(user User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        fmt.Sprint(user.ID),
		Audience:  user.Email,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(cfg.JWT.SESSION_DURATION)).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(cfg.JWT.SECRET))
	return tokenString
}

func generateTestingToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        "1",
		Audience:  "test@test.com",
		Subject:   "testing",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(cfg.JWT.SECRET))
	return tokenString
}

func hash(salt string, data string) string {
	hasher := sha1.New()
	hasher.Write([]byte(salt + data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
