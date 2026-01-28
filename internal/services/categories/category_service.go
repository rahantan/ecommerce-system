package categoryservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	categoryrepositories "ecommerce-system/internal/repositories/categories"
)

type CategoryServiceImpl struct {
	categoryrepositories.CategoryRepositories
}

func NewCategoryService(Category categoryrepositories.CategoryRepositories) CategoryServices {
	return &CategoryServiceImpl{
		CategoryRepositories: Category,
	}
}

func (categoryService *CategoryServiceImpl) handleError(err error) error {
	return exceptions.CheckError(err)
}

func (Category *CategoryServiceImpl) loadCategory(CategoryLoad *models.CategoryModel) *response.ResCategory {
	return &response.ResCategory{
		ID:        CategoryLoad.ID,
		Name:      CategoryLoad.Name,
		CreatedAt: CategoryLoad.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: CategoryLoad.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
func (categoryService *CategoryServiceImpl) GetCategoryById(id int64) (*response.ResCategory, error) {
	result, err := categoryService.CategoryRepositories.GetCategoryById(id)
	if errCheck := categoryService.handleError(err); errCheck != nil {
		return nil, errCheck
	}
	return categoryService.loadCategory(result), nil
}
func (categoryService *CategoryServiceImpl) GetAllCategory() ([]*response.ResCategory, error) {
	result, err := categoryService.CategoryRepositories.GetAllCategory()
	if errCheck := categoryService.handleError(err); errCheck != nil {
		return nil, errCheck
	}

	Categorys := []*response.ResCategory{}
	for _, Category := range result {
		Categorys = append(Categorys, categoryService.loadCategory(Category))
	}

	return Categorys, nil
}

func (categoryService *CategoryServiceImpl) CreateCategory(request *request.ReqCreateCategory) (*response.ResCategory, error) {
	result, err := categoryService.CategoryRepositories.CreateCategory(&models.CategoryModel{
		Name: request.Name,
	})
	if errCheck := categoryService.handleError(err); errCheck != nil {
		return nil, errCheck
	}
	return categoryService.loadCategory(result), nil
}
func (categoryService *CategoryServiceImpl) UpdateCategory(request *request.ReqUpdateCategory) (*response.ResCategory, error) {

	result, err := categoryService.CategoryRepositories.UpdateCategoryById(&models.CategoryModel{
		ID:   request.ID,
		Name: request.Name,
	})
	if errCheck := categoryService.handleError(err); errCheck != nil {
		return nil, errCheck
	}
	return categoryService.loadCategory(result), nil
}
