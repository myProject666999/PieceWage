package service

import (
	"errors"

	"piece-wage/internal/middleware"
	"piece-wage/internal/model"
	"piece-wage/internal/repository"
	"piece-wage/pkg/logger"

	"go.uber.org/zap"
)

type AuthService struct {
	userRepo *repository.UserRepo
}

func NewAuthService() *AuthService {
	return &AuthService{userRepo: repository.NewUserRepo()}
}

func (s *AuthService) Login(req *model.LoginReq) (*model.LoginResp, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		logger.Log.Warn("login failed: user not found", zap.String("username", req.Username))
		return nil, errors.New("用户名或密码错误")
	}

	if !s.userRepo.VerifyPassword(user.Password, req.Password) {
		logger.Log.Warn("login failed: wrong password", zap.String("username", req.Username))
		return nil, errors.New("用户名或密码错误")
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role, user.RealName)
	if err != nil {
		logger.Log.Error("generate token failed", zap.Error(err))
		return nil, errors.New("系统错误")
	}

	logger.Log.Info("login success", zap.String("username", req.Username), zap.Uint64("userId", user.ID))

	return &model.LoginResp{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		RealName: user.RealName,
		Role:     user.Role,
	}, nil
}
