package repository

import (
	"piece-wage/internal/model"
	"piece-wage/pkg/db"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo { return &UserRepo{} }

func (r *UserRepo) GetByUsername(username string) (*model.SysUser, error) {
	var user model.SysUser
	err := db.DB.Where("username = ? AND status = 1", username).First(&user).Error
	return &user, err
}

func (r *UserRepo) GetByID(id uint64) (*model.SysUser, error) {
	var user model.SysUser
	err := db.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepo) List(page, pageSize int, role int) ([]model.SysUser, int64, error) {
	var list []model.SysUser
	var total int64
	q := db.DB.Model(&model.SysUser{}).Where("status = 1")
	if role > 0 {
		q = q.Where("role = ?", role)
	}
	q.Count(&total)
	err := q.Order("id").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *UserRepo) Create(user *model.SysUser) error {
	return db.DB.Create(user).Error
}

func (r *UserRepo) VerifyPassword(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
