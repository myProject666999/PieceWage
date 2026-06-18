package handler

import (
	"strconv"

	"piece-wage/internal/model"
	"piece-wage/internal/service"
	"piece-wage/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProcessHandler struct {
	priceService *service.ProcessPriceService
	stepService  *service.ProcessStepService
	prodService  *service.ProductService
}

func NewProcessHandler() *ProcessHandler {
	return &ProcessHandler{
		priceService: service.NewProcessPriceService(),
		stepService:  service.NewProcessStepService(),
		prodService:  service.NewProductService(),
	}
}

func (h *ProcessHandler) CreatePrice(c *gin.Context) {
	var req model.CreateProcessPriceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailBind(c)
		return
	}
	userID := c.GetUint64("userID")
	price, err := h.priceService.CreatePrice(&req, userID)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, price)
}

func (h *ProcessHandler) GetEffectivePrice(c *gin.Context) {
	processID, _ := strconv.ParseUint(c.Query("processId"), 10, 64)
	gradeLevel := c.DefaultQuery("gradeLevel", "STD")
	date := c.Query("date")
	if processID == 0 || date == "" {
		response.Fail(c, 400, "processId和date不能为空")
		return
	}
	price, err := h.priceService.GetEffectivePrice(processID, gradeLevel, date)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.OK(c, price)
}

func (h *ProcessHandler) ListPricesByProcess(c *gin.Context) {
	processID, _ := strconv.ParseUint(c.Param("processId"), 10, 64)
	list, err := h.priceService.ListByProcess(processID)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, list)
}

func (h *ProcessHandler) ListAllPrices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	list, total, err := h.priceService.ListAll(page, pageSize)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OKPage(c, list, total, page, pageSize)
}

func (h *ProcessHandler) CreateStep(c *gin.Context) {
	var req model.ProcessStepCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailBind(c)
		return
	}
	step, err := h.stepService.Create(&req)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, step)
}

func (h *ProcessHandler) ListSteps(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	list, total, err := h.stepService.List(page, pageSize)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OKPage(c, list, total, page, pageSize)
}

func (h *ProcessHandler) ListStepsByProduct(c *gin.Context) {
	productID, _ := strconv.ParseUint(c.Param("productId"), 10, 64)
	list, err := h.stepService.ListByProduct(productID)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, list)
}

func (h *ProcessHandler) CreateProduct(c *gin.Context) {
	var req model.ProductCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailBind(c)
		return
	}
	p, err := h.prodService.Create(&req)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, p)
}

func (h *ProcessHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	list, total, err := h.prodService.List(page, pageSize)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OKPage(c, list, total, page, pageSize)
}
