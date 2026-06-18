package repository

import (
	"piece-wage/internal/model"
	"piece-wage/pkg/db"
)

type ProcessPriceRepo struct{}

func NewProcessPriceRepo() *ProcessPriceRepo { return &ProcessPriceRepo{} }

func (r *ProcessPriceRepo) Create(price *model.ProcessPrice) error {
	return db.DB.Create(price).Error
}

func (r *ProcessPriceRepo) GetByID(id uint64) (*model.ProcessPrice, error) {
	var p model.ProcessPrice
	err := db.DB.First(&p, id).Error
	return &p, err
}

func (r *ProcessPriceRepo) GetEffectivePrice(processID uint64, gradeLevel, date string) (*model.ProcessPrice, error) {
	var p model.ProcessPrice
	err := db.DB.Where("process_id = ? AND grade_level = ? AND effective_date <= ?", processID, gradeLevel, date).
		Order("effective_date DESC, version_no DESC").
		First(&p).Error
	return &p, err
}

func (r *ProcessPriceRepo) ListByProcess(processID uint64) ([]model.ProcessPrice, error) {
	var list []model.ProcessPrice
	err := db.DB.Where("process_id = ?", processID).
		Order("grade_level, effective_date DESC").
		Find(&list).Error
	return list, err
}

func (r *ProcessPriceRepo) GetCurrentVersionNo(processID uint64, gradeLevel string) (int, error) {
	var maxVer struct{ Max int }
	err := db.DB.Model(&model.ProcessPrice{}).
		Select("COALESCE(MAX(version_no), 0) as max").
		Where("process_id = ? AND grade_level = ?", processID, gradeLevel).
		Scan(&maxVer).Error
	return maxVer.Max, err
}

func (r *ProcessPriceRepo) ExpirePrevious(processID uint64, gradeLevel, effectiveDate string) error {
	return db.DB.Model(&model.ProcessPrice{}).
		Where("process_id = ? AND grade_level = ? AND expiry_date IS NULL AND effective_date < ?", processID, gradeLevel, effectiveDate).
		Update("expiry_date", effectiveDate).Error
}

func (r *ProcessPriceRepo) ListAll(page, pageSize int) ([]model.ProcessPrice, int64, error) {
	var list []model.ProcessPrice
	var total int64
	db.DB.Model(&model.ProcessPrice{}).Count(&total)
	err := db.DB.Preload("Process").Preload("Process.Product").
		Order("effective_date DESC, process_id").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&list).Error
	return list, total, err
}
