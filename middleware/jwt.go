package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func secret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return nil, errors.New("JWT_SECRET is empty")
	}
	return []byte(s), nil
}

func GenerateToken(userID primitive.ObjectID, email, username string, ttl time.Duration) (string, error) {
	key, err := secret()
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	claims := Claims{
		UserID:   userID.Hex(),
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.Hex(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(key)
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		key, err := secret()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return key, nil
		})
		if err != nil || token == nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		uid, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid userId in token"})
			return
		}

		// 注入到 gin context，后续 handler 可用
		c.Set("userID", uid)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)

		c.Next()
	}
}
