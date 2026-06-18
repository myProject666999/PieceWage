package handler

import (
	"strconv"

	"piece-wage/internal/model"
	"piece-wage/internal/service"
	"piece-wage/pkg/response"

	"github.com/gin-gonic/gin"
)

type WageHandler struct {
	wageService *service.WageService
}

func NewWageHandler() *WageHandler {
	return &WageHandler{wageService: service.NewWageService()}
}

func (h *WageHandler) GetMonthlySummary(c *gin.Context) {
	workerID, _ := strconv.ParseUint(c.Param("workerId"), 10, 64)
	month := c.Param("month")
	summary, err := h.wageService.GetMonthlySummary(workerID, month)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.OK(c, summary)
}

func (h *WageHandler) ListSummaries(c *gin.Context) {
	var req model.WageSummaryQueryReq
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
	list, total, err := h.wageService.ListSummaries(&req)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OKPage(c, list, total, req.Page, req.PageSize)
}

func (h *WageHandler) GetWorkerDetails(c *gin.Context) {
	var req model.WorkerWageDetailReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailBind(c)
		return
	}
	list, err := h.wageService.GetWorkerDetails(&req)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, list)
}

func (h *WageHandler) GetWorkerDailyDetails(c *gin.Context) {
	workerID, _ := strconv.ParseUint(c.Param("workerId"), 10, 64)
	wageDate := c.Param("date")
	list, err := h.wageService.GetWorkerDailyDetails(workerID, wageDate)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, list)
}

func (h *WageHandler) GetRealtimeAccumulate(c *gin.Context) {
	workerID, _ := strconv.ParseUint(c.Param("workerId"), 10, 64)
	month := c.Param("month")
	amount, err := h.wageService.GetRealtimeAccumulate(workerID, month)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, gin.H{"workerId": workerID, "month": month, "accumulateAmount": amount})
}

func (h *WageHandler) SettleMonth(c *gin.Context) {
	month := c.Param("month")
	if err := h.wageService.SettleMonth(month); err != nil {
		response.Fail(c, 500, "结算失败: "+err.Error())
		return
	}
	response.OKMsg(c, "月度结算完成")
}
