package service

import (
	"errors"
	"strings"

	"piece-wage/internal/model"
	"piece-wage/internal/repository"
	"piece-wage/pkg/logger"

	"go.uber.org/zap"
)

func isDuplicateEntryErr(err error, keyName string) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Error 1062") && strings.Contains(msg, keyName)
}

type ProcessPriceService struct {
	priceRepo *repository.ProcessPriceRepo
	stepRepo  *repository.ProcessStepRepo
}

func NewProcessPriceService() *ProcessPriceService {
	return &ProcessPriceService{
		priceRepo: repository.NewProcessPriceRepo(),
		stepRepo:  repository.NewProcessStepRepo(),
	}
}

func (s *ProcessPriceService) CreatePrice(req *model.CreateProcessPriceReq, userID uint64) (*model.ProcessPrice, error) {
	step, err := s.stepRepo.GetByID(req.ProcessID)
	if err != nil {
		return nil, errors.New("工序不存在")
	}
	logger.Log.Debug("create price for process",
		zap.Uint64("processId", req.ProcessID),
		zap.String("processName", step.ProcessName),
		zap.Float64("unitPrice", req.UnitPrice),
		zap.String("gradeLevel", req.GradeLevel),
		zap.String("effectiveDate", req.EffectiveDate),
	)

	currentVer, err := s.priceRepo.GetCurrentVersionNo(req.ProcessID, req.GradeLevel)
	if err != nil {
		return nil, err
	}

	if err := s.priceRepo.ExpirePrevious(req.ProcessID, req.GradeLevel, req.EffectiveDate); err != nil {
		logger.Log.Error("expire previous price failed", zap.Error(err))
		return nil, errors.New("更新旧单价失效日期失败")
	}

	price := &model.ProcessPrice{
		ProcessID:     req.ProcessID,
		VersionNo:     currentVer + 1,
		GradeLevel:    req.GradeLevel,
		UnitPrice:     req.UnitPrice,
		EffectiveDate: req.EffectiveDate,
		Remark:        req.Remark,
		CreatedBy:     userID,
	}

	if err := s.priceRepo.Create(price); err != nil {
		logger.Log.Error("create price failed", zap.Error(err))
		return nil, errors.New("创建单价失败")
	}

	logger.Log.Info("price created",
		zap.Uint64("processId", req.ProcessID),
		zap.Int("versionNo", price.VersionNo),
		zap.Float64("unitPrice", price.UnitPrice),
	)

	return price, nil
}

func (s *ProcessPriceService) GetEffectivePrice(processID uint64, gradeLevel, date string) (*model.ProcessPrice, error) {
	price, err := s.priceRepo.GetEffectivePrice(processID, gradeLevel, date)
	if err != nil {
		return nil, errors.New("未找到生效单价")
	}
	logger.Log.Debug("get effective price",
		zap.Uint64("processId", processID),
		zap.String("gradeLevel", gradeLevel),
		zap.String("date", date),
		zap.Float64("unitPrice", price.UnitPrice),
		zap.Int("versionNo", price.VersionNo),
	)
	return price, nil
}

func (s *ProcessPriceService) ListByProcess(processID uint64) ([]model.ProcessPrice, error) {
	return s.priceRepo.ListByProcess(processID)
}

func (s *ProcessPriceService) ListAll(page, pageSize int) ([]model.ProcessPrice, int64, error) {
	return s.priceRepo.ListAll(page, pageSize)
}

type ProcessStepService struct {
	stepRepo    *repository.ProcessStepRepo
	productRepo *repository.ProductRepo
}

func NewProcessStepService() *ProcessStepService {
	return &ProcessStepService{
		stepRepo:    repository.NewProcessStepRepo(),
		productRepo: repository.NewProductRepo(),
	}
}

func (s *ProcessStepService) Create(req *model.ProcessStepCreateReq) (*model.ProcessStep, error) {
	if _, err := s.productRepo.GetByID(req.ProductID); err != nil {
		return nil, errors.New("产品不存在")
	}
	step := &model.ProcessStep{
		ProcessCode: req.ProcessCode,
		ProcessName: req.ProcessName,
		ProductID:   req.ProductID,
		Difficulty:  req.Difficulty,
		Description: req.Description,
		IsShared:    req.IsShared,
	}
	if err := s.stepRepo.Create(step); err != nil {
		if isDuplicateEntryErr(err, "uk_process_code") {
			return nil, errors.New("工序编号已存在")
		}
		logger.Log.Error("create step failed", zap.Error(err))
		return nil, errors.New("创建工序失败")
	}
	return step, nil
}

func (s *ProcessStepService) List(page, pageSize int) ([]model.ProcessStep, int64, error) {
	return s.stepRepo.List(page, pageSize)
}

func (s *ProcessStepService) ListByProduct(productID uint64) ([]model.ProcessStep, error) {
	return s.stepRepo.ListByProduct(productID)
}

type ProductService struct {
	productRepo *repository.ProductRepo
}

func NewProductService() *ProductService {
	return &ProductService{productRepo: repository.NewProductRepo()}
}

func (s *ProductService) Create(req *model.ProductCreateReq) (*model.Product, error) {
	p := &model.Product{
		ProductCode: req.ProductCode,
		ProductName: req.ProductName,
		Spec:        req.Spec,
	}
	if err := s.productRepo.Create(p); err != nil {
		if isDuplicateEntryErr(err, "uk_product_code") {
			return nil, errors.New("产品编号已存在")
		}
		logger.Log.Error("create product failed", zap.Error(err))
		return nil, errors.New("创建产品失败")
	}
	return p, nil
}

func (s *ProductService) List(page, pageSize int) ([]model.Product, int64, error) {
	return s.productRepo.List(page, pageSize)
}
