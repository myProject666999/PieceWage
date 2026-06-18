package service

import (
	"errors"
	"fmt"
	"time"

	"piece-wage/internal/model"
	"piece-wage/internal/repository"
	"piece-wage/pkg/db"
	"piece-wage/pkg/logger"
	"piece-wage/pkg/redis"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductionReportService struct {
	reportRepo    *repository.ProductionReportRepo
	priceRepo     *repository.ProcessPriceRepo
	stepRepo      *repository.ProcessStepRepo
	detailRepo    *repository.WageDetailRepo
	summaryRepo   *repository.WageSummaryRepo
	teamRepo      *repository.TeamRepo
	allocRepo     *repository.AllocationRepo
}

func NewProductionReportService() *ProductionReportService {
	return &ProductionReportService{
		reportRepo:  repository.NewProductionReportRepo(),
		priceRepo:   repository.NewProcessPriceRepo(),
		stepRepo:    repository.NewProcessStepRepo(),
		detailRepo:  repository.NewWageDetailRepo(),
		summaryRepo: repository.NewWageSummaryRepo(),
		teamRepo:    repository.NewTeamRepo(),
		allocRepo:   repository.NewAllocationRepo(),
	}
}

func (s *ProductionReportService) CreateReport(req *model.CreateReportReq, createdBy uint64) (*model.ProductionReport, error) {
	gradeLevel := req.GradeLevel
	if gradeLevel == "" {
		gradeLevel = "STD"
	}

	price, err := s.priceRepo.GetEffectivePrice(req.ProcessID, gradeLevel, req.ReportDate)
	if err != nil {
		return nil, fmt.Errorf("报工日期%s未找到工序%d等级%s的生效单价", req.ReportDate, req.ProcessID, gradeLevel)
	}

	step, err := s.stepRepo.GetByID(req.ProcessID)
	if err != nil {
		return nil, errors.New("工序不存在")
	}

	unitDefect := req.UnitDefect
	if unitDefect == 0 {
		unitDefect = 0.5
	}

	qtyTotal := req.QtyGood + req.QtyDefect
	grossAmount := float64(req.QtyGood) * price.UnitPrice
	defectAmount := float64(req.QtyDefect) * price.UnitPrice * unitDefect
	netAmount := grossAmount - defectAmount

	reportNo := fmt.Sprintf("RPT%s%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)

	report := &model.ProductionReport{
		ReportNo:     reportNo,
		WorkerID:     req.WorkerID,
		TeamID:       req.TeamID,
		ProcessID:    req.ProcessID,
		PriceID:      price.ID,
		UnitPrice:    price.UnitPrice,
		GradeLevel:   gradeLevel,
		ReportDate:   req.ReportDate,
		QtyGood:      req.QtyGood,
		QtyDefect:    req.QtyDefect,
		QtyTotal:     qtyTotal,
		UnitDefect:   unitDefect,
		GrossAmount:  grossAmount,
		DefectAmount: defectAmount,
		NetAmount:    netAmount,
		WorkOrderNo:  req.WorkOrderNo,
		Remark:       req.Remark,
		Status:       1,
		CreatedBy:    createdBy,
	}

	logger.Log.Debug("create report",
		zap.String("reportNo", reportNo),
		zap.Uint64("workerId", req.WorkerID),
		zap.Uint64("processId", req.ProcessID),
		zap.Float64("unitPrice(snap)", price.UnitPrice),
		zap.Int("versionNo", price.VersionNo),
		zap.Int("qtyGood", req.QtyGood),
		zap.Int("qtyDefect", req.QtyDefect),
		zap.Float64("netAmount", netAmount),
	)

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.reportRepo.CreateInTx(tx, report); err != nil {
			return err
		}

		detail := &model.WorkerWageDetail{
			WorkerID:   req.WorkerID,
			ReportID:   &report.ID,
			DetailType: 1,
			ProcessID:  &req.ProcessID,
			WageDate:   req.ReportDate,
			QtyGood:    req.QtyGood,
			QtyDefect:  req.QtyDefect,
			UnitPrice:  price.UnitPrice,
			Amount:     netAmount,
		}
		if err := s.detailRepo.CreateInTx(tx, detail); err != nil {
			return err
		}

		month := req.ReportDate[:7]
		if err := s.summaryRepo.RecalcByWorkerMonth(tx, req.WorkerID, month); err != nil {
			return err
		}

		if step.IsShared == 1 && req.TeamID != nil {
			members, err := s.teamRepo.GetMembers(*req.TeamID)
			if err != nil {
				return err
			}
			if len(members) > 0 {
				perPerson := netAmount / float64(len(members))
				alloc := &model.TeamWageAllocation{
					ReportID:       report.ID,
					TeamID:         *req.TeamID,
					TotalAmount:    netAmount,
					AllocationRule: 1,
					CreatedBy:      createdBy,
				}
				if err := s.allocRepo.CreateInTx(tx, alloc); err != nil {
					return err
				}

				var items []model.TeamWageAllocationItem
				var details []model.WorkerWageDetail
				for _, m := range members {
					items = append(items, model.TeamWageAllocationItem{
						AllocationID: alloc.ID,
						WorkerID:     m.UserID,
						WeightRatio:  1.0 / float64(len(members)),
						AllocatedAmt: perPerson,
					})
					details = append(details, model.WorkerWageDetail{
						WorkerID:     m.UserID,
						AllocationID: &alloc.ID,
						DetailType:   2,
						ProcessID:    &req.ProcessID,
						WageDate:     req.ReportDate,
						Amount:       perPerson,
						Remark:       fmt.Sprintf("班组分配-%s", step.ProcessName),
					})
				}
				if err := s.allocRepo.CreateItemsInTx(tx, items); err != nil {
					return err
				}
				if err := s.detailRepo.CreateBatchInTx(tx, details); err != nil {
					return err
				}

				for _, m := range members {
					if err := s.summaryRepo.RecalcByWorkerMonth(tx, m.UserID, month); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		logger.Log.Error("create report transaction failed", zap.Error(err))
		return nil, errors.New("报工失败")
	}

	s.IncrementRedisCache(req.WorkerID, req.ReportDate[:7], netAmount)

	logger.Log.Info("report created",
		zap.String("reportNo", reportNo),
		zap.Uint64("workerId", req.WorkerID),
		zap.Float64("netAmount", netAmount),
		zap.Bool("isShared", step.IsShared == 1),
	)

	return report, nil
}

func (s *ProductionReportService) GetReport(id uint64) (*model.ProductionReport, error) {
	return s.reportRepo.GetByID(id)
}

func (s *ProductionReportService) ListReports(req *model.ReportQueryReq) ([]model.ProductionReport, int64, error) {
	return s.reportRepo.List(req)
}

func (s *ProductionReportService) VoidReport(id uint64, operatorID uint64) error {
	report, err := s.reportRepo.GetByID(id)
	if err != nil {
		return errors.New("报工单不存在")
	}
	if report.Status == 0 {
		return errors.New("报工单已作废")
	}

	month := report.ReportDate[:7]
	negNetAmount := -report.NetAmount

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.reportRepo.UpdateStatus(id, 0); err != nil {
			return err
		}

		detail := &model.WorkerWageDetail{
			WorkerID:   report.WorkerID,
			ReportID:   &report.ID,
			DetailType: 1,
			ProcessID:  &report.ProcessID,
			WageDate:   report.ReportDate,
			QtyGood:    -report.QtyGood,
			QtyDefect:  -report.QtyDefect,
			UnitPrice:  report.UnitPrice,
			Amount:     negNetAmount,
			Remark:     "作废报工冲销",
		}
		if err := s.detailRepo.CreateInTx(tx, detail); err != nil {
			return err
		}

		if err := s.summaryRepo.RecalcByWorkerMonth(tx, report.WorkerID, month); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.New("作废失败")
	}

	s.IncrementRedisCache(report.WorkerID, month, negNetAmount)

	logger.Log.Info("report voided", zap.Uint64("reportId", id), zap.Uint64("operatorId", operatorID))
	return nil
}

func (s *ProductionReportService) IncrementRedisCache(workerID uint64, month string, amount float64) {
	key := fmt.Sprintf("wage:accumulate:%d:%s", workerID, month)
	val, err := redis.RDB.IncrByFloat(redis.Ctx, key, amount).Result()
	if err != nil {
		logger.Log.Warn("redis increment failed", zap.Error(err))
		return
	}
	redis.RDB.Expire(redis.Ctx, key, 48*time.Hour)
	logger.Log.Debug("redis cache updated", zap.String("key", key), zap.Float64("value", val))
}

func (s *ProductionReportService) GetRedisAccumulate(workerID uint64, month string) (float64, error) {
	key := fmt.Sprintf("wage:accumulate:%d:%s", workerID, month)
	val, err := redis.RDB.Get(redis.Ctx, key).Float64()
	if err != nil {
		if err.Error() == "redis: nil" {
			return 0, nil
		}
		return 0, err
	}
	return val, nil
}
