package handler

import (
	"strconv"

	"piece-wage/internal/model"
	"piece-wage/internal/service"
	"piece-wage/pkg/response"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService  *service.TeamService
	allocService *service.AllocationService
}

func NewTeamHandler() *TeamHandler {
	return &TeamHandler{
		teamService:  service.NewTeamService(),
		allocService: service.NewAllocationService(),
	}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req model.TeamCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailBind(c)
		return
	}
	team, err := h.teamService.Create(&req)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, team)
}

func (h *TeamHandler) ListTeams(c *gin.Context) {
	list, err := h.teamService.List()
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, list)
}

func (h *TeamHandler) GetTeamMembers(c *gin.Context) {
	teamID, _ := strconv.ParseUint(c.Param("teamId"), 10, 64)
	list, err := h.teamService.GetMembers(teamID)
	if err != nil {
		response.Fail(c, 500, "查询失败")
		return
	}
	response.OK(c, list)
}

func (h *TeamHandler) GetAllocation(c *gin.Context) {
	reportID, _ := strconv.ParseUint(c.Param("reportId"), 10, 64)
	alloc, err := h.allocService.GetByReportID(reportID)
	if err != nil {
		response.Fail(c, 404, "分配记录不存在")
		return
	}
	response.OK(c, alloc)
}
