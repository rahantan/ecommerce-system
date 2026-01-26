package categoryrepositories

import "ecommerce-system/internal/models"

type CategoryRepositories interface {
	GetCategoryById(id int64) (*models.CategoryModel, error)
	GetAllCategory() ([]*models.CategoryModel, error)
	UpdateCategoryById(category *models.CategoryModel) (*models.CategoryModel, error)
	CreateCategory(category *models.CategoryModel) (*models.CategoryModel, error)
}
