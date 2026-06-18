package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func OKMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Result{
		Code:    0,
		Message: msg,
	})
}

func Fail(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, Result{
		Code:    -1,
		Message: msg,
	})
}

func FailBind(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Result{
		Code:    -1,
		Message: "参数错误",
	})
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

func OKPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	OK(c, PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
