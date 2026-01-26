package categoryrepositories

import (
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepositoryImpl struct {
	*gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepositories {
	return &CategoryRepositoryImpl{
		DB: db,
	}
}
func (categoryrepo *CategoryRepositoryImpl) GetCategoryById(id int64) (*models.CategoryModel, error) {
	var category models.CategoryModel
	err := categoryrepo.DB.Where("id=?", id).Take(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}
func (categoryrepo *CategoryRepositoryImpl) GetAllCategory() ([]*models.CategoryModel, error) {
	var categories []*models.CategoryModel
	err := categoryrepo.DB.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
func (categoryrepo *CategoryRepositoryImpl) UpdateCategoryById(category *models.CategoryModel) (*models.CategoryModel, error) {
	result := categoryrepo.DB.Model(&models.CategoryModel{}).Where("id=?", category.ID).Updates(category)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, exceptions.ErrNoRowsAffected
	}
	return categoryrepo.GetCategoryById(category.ID)
}
func (categoryrepo *CategoryRepositoryImpl) CreateCategory(category *models.CategoryModel) (*models.CategoryModel, error) {
	err := categoryrepo.DB.Create(category).Error
	if err != nil {
		return nil, err
	}
	return categoryrepo.GetCategoryById(category.ID)
}
