package repository

import (
	"challenge3/models"
	"github.com/jinzhu/gorm"
)

type roleRepoImpl struct {
	DB *gorm.DB
}

func NewRoleRepo(db *gorm.DB) models.RoleRepo {
	return &roleRepoImpl{
		DB: db,
	}
}

func (r *roleRepoImpl) Create(role models.Role) error {
	return r.DB.Create(&role).Error
}

func (r *roleRepoImpl) Find(name string) (models.Role, error) {
	var role models.Role
	result := r.DB.Where("name = ?", name).First(&role)
	if result.Error != nil {
		return models.Role{}, result.Error
	}

	return role, nil
}