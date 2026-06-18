package handler

import (
	"strconv"

	"piece-wage/internal/model"
	"piece-wage/internal/service"
	"piece-wage/pkg/response"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService *service.ProductionReportService
}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{reportService: service.NewProductionReportService()}
}

func (h *ReportHandler) CreateReport(c *gin.Context) {
	var req model.CreateReportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailBind(c)
		return
	}
	userID := c.GetUint64("userID")
	report, err := h.reportService.CreateReport(&req, userID)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, report)
}

func (h *ReportHandler) GetReport(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	report, err := h.reportService.GetReport(id)
	if err != nil {
		response.Fail(c, 404, "报工单不存在")
		return
	}
	response.OK(c, report)
}

func (h *ReportHandler) ListReports(c *gin.Context) {
	var req model.ReportQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailBind(c)
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}
	list, total, err := h.reportService.ListReports(&req)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OKPage(c, list, total, req.Page, req.PageSize)
}

func (h *ReportHandler) VoidReport(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.GetUint64("userID")
	if err := h.reportService.VoidReport(id, userID); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OKMsg(c, "作废成功")
}
