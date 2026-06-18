package repository

import (
	"piece-wage/internal/model"
	"piece-wage/pkg/db"
)

type ProductRepo struct{}

func NewProductRepo() *ProductRepo { return &ProductRepo{} }

func (r *ProductRepo) Create(p *model.Product) error {
	return db.DB.Create(p).Error
}

func (r *ProductRepo) List(page, pageSize int) ([]model.Product, int64, error) {
	var list []model.Product
	var total int64
	db.DB.Model(&model.Product{}).Count(&total)
	err := db.DB.Where("status = 1").Order("id").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *ProductRepo) GetByID(id uint64) (*model.Product, error) {
	var p model.Product
	err := db.DB.First(&p, id).Error
	return &p, err
}

type ProcessStepRepo struct{}

func NewProcessStepRepo() *ProcessStepRepo { return &ProcessStepRepo{} }

func (r *ProcessStepRepo) Create(p *model.ProcessStep) error {
	return db.DB.Create(p).Error
}

func (r *ProcessStepRepo) List(page, pageSize int) ([]model.ProcessStep, int64, error) {
	var list []model.ProcessStep
	var total int64
	db.DB.Model(&model.ProcessStep{}).Count(&total)
	err := db.DB.Preload("Product").Where("status = 1").Order("id").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *ProcessStepRepo) GetByID(id uint64) (*model.ProcessStep, error) {
	var p model.ProcessStep
	err := db.DB.Preload("Product").First(&p, id).Error
	return &p, err
}

func (r *ProcessStepRepo) ListByProduct(productID uint64) ([]model.ProcessStep, error) {
	var list []model.ProcessStep
	err := db.DB.Where("product_id = ? AND status = 1", productID).Order("id").Find(&list).Error
	return list, err
}

func (r *ProcessStepRepo) GetSharedProcesses() ([]model.ProcessStep, error) {
	var list []model.ProcessStep
	err := db.DB.Where("is_shared = 1 AND status = 1").Find(&list).Error
	return list, err
}
