package handler

import (
	"piece-wage/internal/model"
	"piece-wage/internal/service"
	"piece-wage/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{authService: service.NewAuthService()}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailBind(c)
		return
	}
	resp, err := h.authService.Login(&req)
	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}
	response.OK(c, resp)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint64("userID")
	username, _ := c.Get("username")
	realName, _ := c.Get("realName")
	role, _ := c.Get("role")

	response.OK(c, gin.H{
		"userId":   userID,
		"username": username,
		"realName": realName,
		"role":     role,
	})
}
