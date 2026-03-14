package handler

import (
	"errors"
	"net/http"
	"time"
	"week1-postbackend/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"week1-postbackend/middleware"
)

type AuthHandler struct {
	userSvc *service.UserService
}

func NewAuthHandler(userSvc *service.UserService) *AuthHandler {
	return &AuthHandler{userSvc: userSvc}
}

type registerReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	u, err := h.userSvc.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 注意：不要返回 password
	c.JSON(http.StatusCreated, gin.H{
		"id":        u.ID.Hex(),
		"username":  u.Username,
		"email":     u.Email,
		"createdAt": u.Created,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	u, err := h.userSvc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredential) || errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := middleware.GenerateToken(u.ID, u.Email, u.Username, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       u.ID.Hex(),
			"username": u.Username,
			"email":    u.Email,
		},
	})
}
