package service

import (
	"piece-wage/internal/model"
	"piece-wage/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepo
}

func NewUserService() *UserService {
	return &UserService{userRepo: repository.NewUserRepo()}
}

func (s *UserService) List(page, pageSize, role int) ([]model.SysUser, int64, error) {
	return s.userRepo.List(page, pageSize, role)
}

func (s *UserService) GetByID(id uint64) (*model.SysUser, error) {
	return s.userRepo.GetByID(id)
}
