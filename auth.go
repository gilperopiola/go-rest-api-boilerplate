package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gilperopiola/go-rest-api-boilerplate/utils"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	Email          string
	Password       string
	RepeatPassword string
	Code           string
}

func validateToken(requiredRoles ...Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")

		if len(tokenString) < 40 {
			c.JSON(http.StatusUnauthorized, "authentication error 1")
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.SECRET), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, "authentication error 2")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			if hasRequiredRoles(utils.ToInt(claims.Id), requiredRoles) {
				c.Set("ID", claims.Id)
				c.Set("Email", claims.Audience)
			} else {
				c.JSON(http.StatusUnauthorized, "authentication error 3")
				c.Abort()
			}
		} else {
			c.JSON(http.StatusUnauthorized, "authentication error 4")
			c.Abort()
		}
	}
}

func hasRequiredRoles(IDUser int, roles []Role) bool {
	rolesStrings := make([]string, 0)
	for _, role := range roles {
		rolesStrings = append(rolesStrings, utils.ToString(int(role)))
	}
	rolesString := strings.Join(rolesStrings, ",")

	var ID int
	err := db.DB.QueryRow(fmt.Sprintf(`SELECT id FROM users_roles WHERE idUser = %d AND idRole IN (%s)`, IDUser, rolesString)).Scan(&ID)
	if err != nil {
		return false
	}

	return true
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

func generateTestingToken(roles ...Role) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        "9999",
		Audience:  "test@test.com",
		Subject:   "testing",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
	})

	for _, role := range roles {
		db.DB.Exec(`INSERT INTO users_roles (idUser, idRole) VALUES (?, ?)`, 9999, role)
	}

	tokenString, _ := token.SignedString([]byte(cfg.JWT.SECRET))
	return tokenString
}

func hash(salt string, data string) string {
	hasher := sha1.New()
	hasher.Write([]byte(salt + data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
