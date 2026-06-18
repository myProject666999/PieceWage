package handler

import (
	"piece-wage/internal/service"
	"piece-wage/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{userService: service.NewUserService()}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	role, _ := strconv.Atoi(c.DefaultQuery("role", "0"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	list, total, err := h.userService.List(page, pageSize, role)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OKPage(c, list, total, page, pageSize)
}
