package services

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/domain"
	"github.com/hutamatr/go-blog-api/model/web"
	repositoriesC "github.com/hutamatr/go-blog-api/repositories/categories"
)

type CategoryServiceImpl struct {
	repository repositoriesC.CategoryRepository
	db         *sql.DB
	validator  *validator.Validate
}

func NewCategoryService(categoryRepository repositoriesC.CategoryRepository, db *sql.DB, validator *validator.Validate) CategoryService {
	return &CategoryServiceImpl{
		repository: categoryRepository,
		db:         db,
		validator:  validator,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	newCategory := domain.Category{
		Name: request.Name,
	}

	createdCategory := service.repository.Save(ctx, tx, newCategory)

	return web.CategoryResponse(createdCategory)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []web.CategoryResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	categories := service.repository.FindAll(ctx, tx)

	var categoriesData []web.CategoryResponse

	for _, category := range categories {
		categoriesData = append(categoriesData, web.CategoryResponse(category))
	}

	return categoriesData
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) web.CategoryResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	category := service.repository.FindById(ctx, tx, categoryId)

	return web.CategoryResponse(category)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	categoryData := service.repository.FindById(ctx, tx, request.Id)

	categoryData.Name = request.Name

	updatedCategory := service.repository.Update(ctx, tx, categoryData)

	return web.CategoryResponse(updatedCategory)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	categoryData := service.repository.FindById(ctx, tx, categoryId)

	service.repository.Delete(ctx, tx, categoryData.Id)
}
