package categoryrepositories

import (
	"ecommerce-system/internal/models"

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

func (categoryRepo *CategoryRepositoryImpl) checkErrMysql(err error) error {
	if models.IsInternalErrMysql(err) {
		return err
	}
	return models.ErrCategoryNotFound
}

func (categoryRepo *CategoryRepositoryImpl) checkNotFoundForUpdate(categoryID int64) bool {

	var count int64
	if err := categoryRepo.DB.Model(&models.CategoryModel{}).Where("id = ?", categoryID).Count(&count).Error; err != nil {
		return false
	}
	return count == 0
}

func (categoryRepo *CategoryRepositoryImpl) GetCategoryById(categoryID int64) (*models.CategoryModel, error) {
	var category models.CategoryModel

	if err := categoryRepo.DB.Where("id=?", categoryID).Take(&category).Error; err != nil {
		return nil, categoryRepo.checkErrMysql(err)
	}

	return &category, nil
}

func (categoryRepo *CategoryRepositoryImpl) GetAllCategory() ([]*models.CategoryModel, error) {

	var categories []*models.CategoryModel
	if err := categoryRepo.DB.Find(&categories).Error; err != nil {
		return nil, categoryRepo.checkErrMysql(err)
	}

	return categories, nil
}

func (categoryRepo *CategoryRepositoryImpl) UpdateCategoryById(category *models.CategoryModel) (*models.CategoryModel, error) {

	if categoryRepo.checkNotFoundForUpdate(category.ID) {
		return nil, models.ErrCategoryNotFound
	}

	result := categoryRepo.DB.Model(&models.CategoryModel{}).Where("id=?", category.ID).Updates(category)
	if result.Error != nil {
		return nil, categoryRepo.checkErrMysql(result.Error)
	}
	return categoryRepo.GetCategoryById(category.ID)
}

func (categoryRepo *CategoryRepositoryImpl) CreateCategory(category *models.CategoryModel) (*models.CategoryModel, error) {

	err := categoryRepo.DB.Create(category).Error
	if err != nil {
		return nil, categoryRepo.checkErrMysql(err)
	}

	return categoryRepo.GetCategoryById(category.ID)
}
