package repository

import (
	"piece-wage/internal/model"
	"piece-wage/pkg/db"

	"gorm.io/gorm"
)

type ProductionReportRepo struct{}

func NewProductionReportRepo() *ProductionReportRepo { return &ProductionReportRepo{} }

func (r *ProductionReportRepo) Create(report *model.ProductionReport) error {
	return db.DB.Create(report).Error
}

func (r *ProductionReportRepo) GetByID(id uint64) (*model.ProductionReport, error) {
	var report model.ProductionReport
	err := db.DB.Preload("Worker").Preload("Process").Preload("Process.Product").Preload("Price").
		First(&report, id).Error
	return &report, err
}

func (r *ProductionReportRepo) List(req *model.ReportQueryReq) ([]model.ProductionReport, int64, error) {
	var list []model.ProductionReport
	var total int64
	q := db.DB.Model(&model.ProductionReport{})

	if req.WorkerID > 0 {
		q = q.Where("worker_id = ?", req.WorkerID)
	}
	if req.TeamID > 0 {
		q = q.Where("team_id = ?", req.TeamID)
	}
	if req.ProcessID > 0 {
		q = q.Where("process_id = ?", req.ProcessID)
	}
	if req.StartDate != "" {
		q = q.Where("report_date >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		q = q.Where("report_date <= ?", req.EndDate)
	}
	if req.Status > 0 {
		q = q.Where("status = ?", req.Status)
	}

	q.Count(&total)
	err := q.Preload("Worker").Preload("Process").Preload("Process.Product").
		Order("report_date DESC, id DESC").
		Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).
		Find(&list).Error
	return list, total, err
}

func (r *ProductionReportRepo) UpdateStatus(id uint64, status int) error {
	return db.DB.Model(&model.ProductionReport{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ProductionReportRepo) SumNetAmountByWorkerMonth(workerID uint64, month string) (float64, error) {
	var result struct{ Total float64 }
	err := db.DB.Model(&model.ProductionReport{}).
		Select("COALESCE(SUM(net_amount), 0) as total").
		Where("worker_id = ? AND DATE_FORMAT(report_date, '%Y-%m') = ? AND status = 1", workerID, month).
		Scan(&result).Error
	return result.Total, err
}

func (r *ProductionReportRepo) FindByWorkerAndMonth(workerID uint64, month string) ([]model.ProductionReport, error) {
	var list []model.ProductionReport
	err := db.DB.Preload("Process").Preload("Process.Product").
		Where("worker_id = ? AND DATE_FORMAT(report_date, '%Y-%m') = ? AND status = 1", workerID, month).
		Order("report_date DESC").
		Find(&list).Error
	return list, err
}

func (r *ProductionReportRepo) ExistsByReportNo(reportNo string) (bool, error) {
	var count int64
	err := db.DB.Model(&model.ProductionReport{}).Where("report_no = ?", reportNo).Count(&count).Error
	return count > 0, err
}

func (r *ProductionReportRepo) CreateInTx(tx *gorm.DB, report *model.ProductionReport) error {
	return tx.Create(report).Error
}
