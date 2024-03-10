package go_blog

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/go-blog-api/cmd/go_blog/helper"
)

type CategoryService interface {
	Create(ctx context.Context, category CategoryCreateRequest) CategoryResponse
	FindAll(ctx context.Context) []CategoryResponse
	FindById(ctx context.Context, categoryId int) CategoryResponse
	Update(ctx context.Context, category CategoryUpdateRequest) CategoryResponse
	Delete(ctx context.Context, categoryId int)
}

type CategoryServiceImpl struct {
	repository CategoryRepository
	db         *sql.DB
	validator  *validator.Validate
}

func NewCategoryService(categoryRepository CategoryRepository, db *sql.DB, validator *validator.Validate) CategoryService {
	return &CategoryServiceImpl{
		repository: categoryRepository,
		db:         db,
		validator:  validator,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, category CategoryCreateRequest) CategoryResponse {
	err := service.validator.Struct(category)
	helper.PanicError(err)

	tx, err := service.db.Begin()
	helper.PanicError(err)
	defer helper.TxRollbackCommit(tx)

	newCategory := Category{
		Name: category.Name,
	}

	createdCategory := service.repository.Save(ctx, tx, newCategory)

	return CategoryResponse(createdCategory)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []CategoryResponse {
	tx, err := service.db.Begin()
	helper.PanicError(err)
	defer helper.TxRollbackCommit(tx)

	categories := service.repository.FindAll(ctx, tx)

	var categoriesData []CategoryResponse

	for _, category := range categories {
		categoriesData = append(categoriesData, CategoryResponse(category))
	}

	return categoriesData
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) CategoryResponse {
	tx, err := service.db.Begin()
	helper.PanicError(err)
	defer helper.TxRollbackCommit(tx)

	category, err := service.repository.FindById(ctx, tx, categoryId)
	helper.PanicError(err)

	return CategoryResponse(category)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, category CategoryUpdateRequest) CategoryResponse {
	tx, err := service.db.Begin()
	helper.PanicError(err)
	defer helper.TxRollbackCommit(tx)

	categoryData, err := service.repository.FindById(ctx, tx, category.Id)
	helper.PanicError(err)

	categoryData.Name = category.Name

	updatedCategory := service.repository.Update(ctx, tx, categoryData)

	return CategoryResponse(updatedCategory)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) {
	tx, err := service.db.Begin()
	helper.PanicError(err)
	defer helper.TxRollbackCommit(tx)

	categoryData, err := service.repository.FindById(ctx, tx, categoryId)
	helper.PanicError(err)

	service.repository.Delete(ctx, tx, categoryData.Id)
}
