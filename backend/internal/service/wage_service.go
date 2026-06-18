package service

import (
	"errors"

	"piece-wage/internal/model"
	"piece-wage/internal/repository"
	"piece-wage/pkg/db"
	"piece-wage/pkg/logger"
	"piece-wage/pkg/redis"

	"fmt"

	"go.uber.org/zap"
)

type WageService struct {
	summaryRepo *repository.WageSummaryRepo
	detailRepo  *repository.WageDetailRepo
	reportRepo  *repository.ProductionReportRepo
}

func NewWageService() *WageService {
	return &WageService{
		summaryRepo: repository.NewWageSummaryRepo(),
		detailRepo:  repository.NewWageDetailRepo(),
		reportRepo:  repository.NewProductionReportRepo(),
	}
}

func (s *WageService) GetMonthlySummary(workerID uint64, month string) (*model.WorkerWageSummary, error) {
	summary, err := s.summaryRepo.GetByWorkerMonth(workerID, month)
	if err != nil {
		return nil, errors.New("未找到汇总记录")
	}
	return summary, nil
}

func (s *WageService) ListSummaries(req *model.WageSummaryQueryReq) ([]model.WorkerWageSummary, int64, error) {
	return s.summaryRepo.List(req)
}

func (s *WageService) GetWorkerDetails(req *model.WorkerWageDetailReq) ([]model.WorkerWageDetail, error) {
	return s.detailRepo.FindByWorkerAndDateRange(req.WorkerID, req.StartDate, req.EndDate)
}

func (s *WageService) GetWorkerDailyDetails(workerID uint64, wageDate string) ([]model.WorkerWageDetail, error) {
	return s.detailRepo.FindByWorkerAndDate(workerID, wageDate)
}

func (s *WageService) GetRealtimeAccumulate(workerID uint64, month string) (float64, error) {
	cached, err := reportServiceSingleton.GetRedisAccumulate(workerID, month)
	if err == nil && cached > 0 {
		logger.Log.Debug("realtime accumulate from redis", zap.Uint64("workerId", workerID), zap.Float64("amount", cached))
		return cached, nil
	}

	summary, err := s.summaryRepo.GetByWorkerMonth(workerID, month)
	if err != nil {
		return 0, nil
	}

	key := fmt.Sprintf("wage:accumulate:%d:%s", workerID, month)
	redis.RDB.Set(redis.Ctx, key, summary.NetAmount, 0)

	return summary.NetAmount, nil
}

var reportServiceSingleton = NewProductionReportService()

func (s *WageService) SettleMonth(month string) error {
	logger.Log.Info("start settling month", zap.String("month", month))
	return db.DB.Exec(`
		INSERT INTO worker_wage_summary (worker_id, summary_month, total_qty_good, total_qty_defect,
			gross_amount, defect_amount, allocation_amt, adjust_amount, net_amount, calc_status, last_calc_time)
		SELECT
			pr.worker_id,
			?,
			COALESCE(SUM(pr.qty_good), 0),
			COALESCE(SUM(pr.qty_defect), 0),
			COALESCE(SUM(pr.gross_amount), 0),
			COALESCE(SUM(pr.defect_amount), 0),
			0,
			0,
			COALESCE(SUM(pr.net_amount), 0),
			1,
			NOW()
		FROM production_report pr
		WHERE DATE_FORMAT(pr.report_date, '%Y-%m') = ? AND pr.status = 1
		GROUP BY pr.worker_id
		ON DUPLICATE KEY UPDATE
			total_qty_good = VALUES(total_qty_good),
			total_qty_defect = VALUES(total_qty_defect),
			gross_amount = VALUES(gross_amount),
			defect_amount = VALUES(defect_amount),
			net_amount = VALUES(net_amount),
			calc_status = 1,
			last_calc_time = NOW()
	`, month, month).Error
}
