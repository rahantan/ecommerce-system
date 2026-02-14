package usecase

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"

	"gorm.io/gorm"
)

type CategoryUseCaseImpl struct {
	domain.CategoryRepository
	*gorm.DB
}

func NewCategoryUseCase(Category domain.CategoryRepository, db *gorm.DB) domain.CategoryUseCase {
	return &CategoryUseCaseImpl{
		CategoryRepository: Category,
		DB:                 db,
	}
}

func (Category *CategoryUseCaseImpl) loadCategory(CategoryLoad *model.CategoryModel) *response.ResCategory {
	return &response.ResCategory{
		ID:        CategoryLoad.ID,
		Name:      CategoryLoad.Name,
		CreatedAt: CategoryLoad.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: CategoryLoad.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
func (categoryUC *CategoryUseCaseImpl) GetCategoryById(categoryID int64) (*response.ResCategory, error) {

	result, err := categoryUC.CategoryRepository.GetCategoryById(categoryUC.DB, categoryID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return categoryUC.loadCategory(result), nil
}
func (categoryUC *CategoryUseCaseImpl) GetAllCategory() ([]*response.ResCategory, error) {

	result, err := categoryUC.CategoryRepository.GetAllCategory(categoryUC.DB)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	Categorys := []*response.ResCategory{}
	for _, Category := range result {
		Categorys = append(Categorys, categoryUC.loadCategory(Category))
	}

	return Categorys, nil
}

func (categoryUC *CategoryUseCaseImpl) CreateCategory(request *request.ReqCreateCategory) (*response.ResCategory, error) {
	result, err := categoryUC.CategoryRepository.CreateCategory(categoryUC.DB, &model.CategoryModel{
		Name: request.Name,
	})

	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return categoryUC.loadCategory(result), nil
}
func (categoryUC *CategoryUseCaseImpl) UpdateCategory(request *request.ReqUpdateCategory, categoryID int64) (*response.ResCategory, error) {

	result, err := categoryUC.CategoryRepository.UpdateCategoryById(categoryUC.DB, &model.CategoryModel{
		ID:   categoryID,
		Name: request.Name,
	})

	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return categoryUC.loadCategory(result), nil
}
