package categoryservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/models"
	categoryrepositories "ecommerce-system/internal/repositories/categories"
	"ecommerce-system/internal/utils"

	"gorm.io/gorm"
)

type CategoryServiceImpl struct {
	categoryrepositories.CategoryRepositories
	*gorm.DB
}

func NewCategoryService(Category categoryrepositories.CategoryRepositories, db *gorm.DB) CategoryServices {
	return &CategoryServiceImpl{
		CategoryRepositories: Category,
		DB:                   db,
	}
}

func (Category *CategoryServiceImpl) loadCategory(CategoryLoad *models.CategoryModel) *response.ResCategory {
	return &response.ResCategory{
		ID:        CategoryLoad.ID,
		Name:      CategoryLoad.Name,
		CreatedAt: CategoryLoad.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: CategoryLoad.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
func (categoryService *CategoryServiceImpl) GetCategoryById(categoryID int64) (*response.ResCategory, error) {

	result, err := categoryService.CategoryRepositories.GetCategoryById(categoryService.DB, categoryID)
	if err != nil {
		return nil, utils.MappingError(err)
	}

	return categoryService.loadCategory(result), nil
}
func (categoryService *CategoryServiceImpl) GetAllCategory() ([]*response.ResCategory, error) {

	result, err := categoryService.CategoryRepositories.GetAllCategory(categoryService.DB)
	if err != nil {
		return nil, utils.MappingError(err)
	}

	Categorys := []*response.ResCategory{}
	for _, Category := range result {
		Categorys = append(Categorys, categoryService.loadCategory(Category))
	}

	return Categorys, nil
}

func (categoryService *CategoryServiceImpl) CreateCategory(request *request.ReqCreateCategory) (*response.ResCategory, error) {
	result, err := categoryService.CategoryRepositories.CreateCategory(categoryService.DB, &models.CategoryModel{
		Name: request.Name,
	})

	if err != nil {
		return nil, utils.MappingError(err)
	}

	return categoryService.loadCategory(result), nil
}
func (categoryService *CategoryServiceImpl) UpdateCategory(request *request.ReqUpdateCategory, categoryID int64) (*response.ResCategory, error) {

	result, err := categoryService.CategoryRepositories.UpdateCategoryById(categoryService.DB, &models.CategoryModel{
		ID:   categoryID,
		Name: request.Name,
	})

	if err != nil {
		return nil, utils.MappingError(err)
	}

	return categoryService.loadCategory(result), nil
}
