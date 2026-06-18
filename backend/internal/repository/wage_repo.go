package repository

import (
	"piece-wage/internal/model"
	"piece-wage/pkg/db"

	"gorm.io/gorm"
)

type TeamRepo struct{}

func NewTeamRepo() *TeamRepo { return &TeamRepo{} }

func (r *TeamRepo) Create(team *model.WorkTeam) error {
	return db.DB.Create(team).Error
}

func (r *TeamRepo) GetByID(id uint64) (*model.WorkTeam, error) {
	var t model.WorkTeam
	err := db.DB.First(&t, id).Error
	return &t, err
}

func (r *TeamRepo) List() ([]model.WorkTeam, error) {
	var list []model.WorkTeam
	err := db.DB.Where("status = 1").Order("id").Find(&list).Error
	return list, err
}

func (r *TeamRepo) GetMembers(teamID uint64) ([]model.TeamMember, error) {
	var list []model.TeamMember
	err := db.DB.Where("team_id = ? AND leave_date IS NULL", teamID).Find(&list).Error
	return list, err
}

type AllocationRepo struct{}

func NewAllocationRepo() *AllocationRepo { return &AllocationRepo{} }

func (r *AllocationRepo) CreateInTx(tx *gorm.DB, alloc *model.TeamWageAllocation) error {
	return tx.Create(alloc).Error
}

func (r *AllocationRepo) CreateItemsInTx(tx *gorm.DB, items []model.TeamWageAllocationItem) error {
	return tx.Create(&items).Error
}

func (r *AllocationRepo) GetByReportID(reportID uint64) (*model.TeamWageAllocation, error) {
	var a model.TeamWageAllocation
	err := db.DB.Preload("Items").Preload("Items.Worker").
		Where("report_id = ?", reportID).First(&a).Error
	return &a, err
}

type WageDetailRepo struct{}

func NewWageDetailRepo() *WageDetailRepo { return &WageDetailRepo{} }

func (r *WageDetailRepo) CreateInTx(tx *gorm.DB, detail *model.WorkerWageDetail) error {
	return tx.Create(detail).Error
}

func (r *WageDetailRepo) CreateBatchInTx(tx *gorm.DB, details []model.WorkerWageDetail) error {
	return tx.Create(&details).Error
}

func (r *WageDetailRepo) FindByWorkerAndDateRange(workerID uint64, startDate, endDate string) ([]model.WorkerWageDetail, error) {
	var list []model.WorkerWageDetail
	q := db.DB.Where("worker_id = ?", workerID)
	if startDate != "" {
		q = q.Where("wage_date >= ?", startDate)
	}
	if endDate != "" {
		q = q.Where("wage_date <= ?", endDate)
	}
	err := q.Order("wage_date DESC, id DESC").Find(&list).Error
	return list, err
}

func (r *WageDetailRepo) FindByWorkerAndDate(workerID uint64, wageDate string) ([]model.WorkerWageDetail, error) {
	var list []model.WorkerWageDetail
	err := db.DB.Where("worker_id = ? AND wage_date = ?", workerID, wageDate).
		Order("id DESC").Find(&list).Error
	return list, err
}

type WageSummaryRepo struct{}

func NewWageSummaryRepo() *WageSummaryRepo { return &WageSummaryRepo{} }

func (r *WageSummaryRepo) Upsert(summary *model.WorkerWageSummary) error {
	return db.DB.Save(summary).Error
}

func (r *WageSummaryRepo) GetByWorkerMonth(workerID uint64, month string) (*model.WorkerWageSummary, error) {
	var s model.WorkerWageSummary
	err := db.DB.Preload("Worker").Where("worker_id = ? AND summary_month = ?", workerID, month).First(&s).Error
	return &s, err
}

func (r *WageSummaryRepo) List(req *model.WageSummaryQueryReq) ([]model.WorkerWageSummary, int64, error) {
	var list []model.WorkerWageSummary
	var total int64
	q := db.DB.Model(&model.WorkerWageSummary{})

	if req.WorkerID > 0 {
		q = q.Where("worker_id = ?", req.WorkerID)
	}
	if req.SummaryMonth != "" {
		q = q.Where("summary_month = ?", req.SummaryMonth)
	}

	q.Count(&total)
	err := q.Preload("Worker").
		Order("summary_month DESC, worker_id").
		Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).
		Find(&list).Error
	return list, total, err
}

func (r *WageSummaryRepo) RecalcByWorkerMonth(tx *gorm.DB, workerID uint64, month string) error {
	var directResult struct {
		GrossAmt    float64
		DefectAmt   float64
		QtyGood     int
		QtyDefect   int
	}
	err := tx.Model(&model.ProductionReport{}).
		Select("COALESCE(SUM(gross_amount),0) as gross_amt, COALESCE(SUM(defect_amount),0) as defect_amt, COALESCE(SUM(qty_good),0) as qty_good, COALESCE(SUM(qty_defect),0) as qty_defect").
		Where("worker_id = ? AND DATE_FORMAT(report_date,'%Y-%m') = ? AND status = 1", workerID, month).
		Scan(&directResult).Error
	if err != nil {
		return err
	}

	var allocResult struct{ Total float64 }
	err = tx.Model(&model.TeamWageAllocationItem{}).
		Select("COALESCE(SUM(twai.allocated_amt),0) as total").
		Joins("JOIN team_wage_allocation twa ON twa.id = twai.allocation_id").
		Joins("JOIN production_report pr ON pr.id = twa.report_id").
		Where("twai.worker_id = ? AND DATE_FORMAT(pr.report_date,'%Y-%m') = ? AND pr.status = 1", workerID, month).
		Scan(&allocResult).Error
	if err != nil {
		return err
	}

	netAmount := directResult.GrossAmt - directResult.DefectAmt + allocResult.Total

	summary := &model.WorkerWageSummary{
		WorkerID:       workerID,
		SummaryMonth:   month,
		TotalQtyGood:   directResult.QtyGood,
		TotalQtyDefect: directResult.QtyDefect,
		GrossAmount:    directResult.GrossAmt,
		DefectAmount:   directResult.DefectAmt,
		AllocationAmt:  allocResult.Total,
		NetAmount:      netAmount,
	}

	return tx.Save(summary).Error
}
