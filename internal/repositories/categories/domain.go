package categoryrepositories

import (
	"ecommerce-system/internal/models"

	"gorm.io/gorm"
)

type CategoryRepositories interface {
	GetCategoryById(db *gorm.DB, categoryID int64) (*models.CategoryModel, error)
	GetAllCategory(db *gorm.DB) ([]*models.CategoryModel, error)
	UpdateCategoryById(db *gorm.DB, category *models.CategoryModel) (*models.CategoryModel, error)
	CreateCategory(db *gorm.DB, category *models.CategoryModel) (*models.CategoryModel, error)
}
