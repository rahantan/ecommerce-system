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

func (Category *CategoryServiceImpl) loadCategory(CategoryLoad *models.CategoryModel) *response.ResCategory {
	return &response.ResCategory{
		ID:        CategoryLoad.ID,
		Name:      CategoryLoad.Name,
		CreatedAt: CategoryLoad.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: CategoryLoad.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
func (CategoryService *CategoryServiceImpl) GetCategoryById(id int64) (*response.ResCategory, error) {
	result, err := CategoryService.CategoryRepositories.GetCategoryById(id)
	if errCheck := exceptions.CheckError(err); errCheck != nil {
		return nil, errCheck
	}
	return CategoryService.loadCategory(result), nil
}
func (CategoryService *CategoryServiceImpl) GetAllCategory() ([]*response.ResCategory, error) {
	result, err := CategoryService.CategoryRepositories.GetAllCategory()
	if errCheck := exceptions.CheckError(err); errCheck != nil {
		return nil, errCheck
	}

	Categorys := []*response.ResCategory{}
	for _, Category := range result {
		Categorys = append(Categorys, CategoryService.loadCategory(Category))
	}

	return Categorys, nil
}

func (CategoryService *CategoryServiceImpl) CreateCategory(request *request.ReqCreateCategory) (*response.ResCategory, error) {
	result, err := CategoryService.CategoryRepositories.CreateCategory(&models.CategoryModel{
		Name: request.Name,
	})
	if errCheck := exceptions.CheckError(err); errCheck != nil {
		return nil, errCheck
	}
	return CategoryService.loadCategory(result), nil
}
func (CategoryService *CategoryServiceImpl) UpdateCategory(request *request.ReqUpdateCategory) (*response.ResCategory, error) {

	result, err := CategoryService.CategoryRepositories.UpdateCategoryById(&models.CategoryModel{
		ID:   request.ID,
		Name: request.Name,
	})
	if errCheck := exceptions.CheckError(err); errCheck != nil {
		return nil, errCheck
	}
	return CategoryService.loadCategory(result), nil
}
