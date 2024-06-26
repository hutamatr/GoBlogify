package category

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
)

type CategoryService interface {
	Create(ctx context.Context, request CategoryCreateRequest, isAdmin bool) CategoryResponse
	FindAll(ctx context.Context, limit, offset int) ([]CategoryResponse, int)
	FindById(ctx context.Context, categoryId int) CategoryResponse
	Update(ctx context.Context, request CategoryUpdateRequest, isAdmin bool) CategoryResponse
	Delete(ctx context.Context, categoryId int, isAdmin bool)
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

func (service *CategoryServiceImpl) Create(ctx context.Context, request CategoryCreateRequest, isAdmin bool) CategoryResponse {
	if !isAdmin {
		panic(exception.NewBadRequestError("only admin can create category"))
	}

	err := service.validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	newCategory := Category{
		Name: request.Name,
	}

	createdCategory := service.repository.Save(ctx, tx, newCategory)

	return ToCategoryResponse(createdCategory)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context, limit, offset int) ([]CategoryResponse, int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	categories := service.repository.FindAll(ctx, tx, limit, offset)
	countCategories := service.repository.CountCategories(ctx, tx)

	var categoriesData []CategoryResponse

	if len(categories) == 0 {
		panic(exception.NewNotFoundError("categories not found"))
	}

	for _, category := range categories {
		categoriesData = append(categoriesData, ToCategoryResponse(category))
	}

	return categoriesData, countCategories
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) CategoryResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	category := service.repository.FindById(ctx, tx, categoryId)

	return ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request CategoryUpdateRequest, isAdmin bool) CategoryResponse {
	if !isAdmin {
		panic(exception.NewBadRequestError("only admin can update category"))
	}

	err := service.validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	categoryData := service.repository.FindById(ctx, tx, request.Id)

	categoryData.Name = request.Name

	updatedCategory := service.repository.Update(ctx, tx, categoryData)

	return ToCategoryResponse(updatedCategory)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int, isAdmin bool) {
	if !isAdmin {
		panic(exception.NewBadRequestError("only admin can delete category"))
	}

	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	service.repository.Delete(ctx, tx, categoryId)
}
