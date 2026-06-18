package service

import (
	"piece-wage/internal/model"
	"piece-wage/internal/repository"
)

type TeamService struct {
	teamRepo *repository.TeamRepo
}

func NewTeamService() *TeamService {
	return &TeamService{teamRepo: repository.NewTeamRepo()}
}

func (s *TeamService) Create(req *model.TeamCreateReq) (*model.WorkTeam, error) {
	team := &model.WorkTeam{
		TeamName: req.TeamName,
		TeamCode: req.TeamCode,
		LeaderID: req.LeaderID,
	}
	if err := s.teamRepo.Create(team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *TeamService) List() ([]model.WorkTeam, error) {
	return s.teamRepo.List()
}

func (s *TeamService) GetMembers(teamID uint64) ([]model.TeamMember, error) {
	return s.teamRepo.GetMembers(teamID)
}

type AllocationService struct {
	allocRepo *repository.AllocationRepo
}

func NewAllocationService() *AllocationService {
	return &AllocationService{allocRepo: repository.NewAllocationRepo()}
}

func (s *AllocationService) GetByReportID(reportID uint64) (*model.TeamWageAllocation, error) {
	return s.allocRepo.GetByReportID(reportID)
}
