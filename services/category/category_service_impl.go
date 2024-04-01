package services

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/domain"
	"github.com/hutamatr/GoBlogify/model/web"
	repositoriesCategory "github.com/hutamatr/GoBlogify/repositories/category"
)

type CategoryServiceImpl struct {
	repository repositoriesCategory.CategoryRepository
	db         *sql.DB
	validator  *validator.Validate
}

func NewCategoryService(categoryRepository repositoriesCategory.CategoryRepository, db *sql.DB, validator *validator.Validate) CategoryService {
	return &CategoryServiceImpl{
		repository: categoryRepository,
		db:         db,
		validator:  validator,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest, isAdmin bool) web.CategoryResponse {
	if !isAdmin {
		panic(exception.NewBadRequestError("only admin can create category"))
	}

	err := service.validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	newCategory := domain.Category{
		Name: request.Name,
	}

	createdCategory := service.repository.Save(ctx, tx, newCategory)

	return web.ToCategoryResponse(createdCategory)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []web.CategoryResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	categories := service.repository.FindAll(ctx, tx)

	var categoriesData []web.CategoryResponse

	if len(categories) == 0 {
		return categoriesData
	}

	for _, category := range categories {
		categoriesData = append(categoriesData, web.ToCategoryResponse(category))
	}

	return categoriesData
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) web.CategoryResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	category := service.repository.FindById(ctx, tx, categoryId)

	return web.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest, isAdmin bool) web.CategoryResponse {
	if !isAdmin {
		panic(exception.NewBadRequestError("only admin can update category"))
	}

	err := service.validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	categoryData := service.repository.FindById(ctx, tx, request.Id)

	categoryData.Name = request.Name

	updatedCategory := service.repository.Update(ctx, tx, categoryData)

	return web.ToCategoryResponse(updatedCategory)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int, isAdmin bool) {
	if !isAdmin {
		panic(exception.NewBadRequestError("only admin can delete category"))
	}

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	service.repository.Delete(ctx, tx, categoryId)
}
