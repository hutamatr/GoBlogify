package category

import (
	"context"
	"database/sql"

	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
)

type CategoryRepository interface {
	Save(ctx context.Context, tx *sql.Tx, category Category) Category
	FindAll(ctx context.Context, tx *sql.Tx, limit, offset int) []Category
	FindById(ctx context.Context, tx *sql.Tx, categoryId int) Category
	Update(ctx context.Context, tx *sql.Tx, category Category) Category
	Delete(ctx context.Context, tx *sql.Tx, categoryId int)
	CountCategories(ctx context.Context, tx *sql.Tx) int
}

type CategoryRepositoryImpl struct {
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
}

func (repository *CategoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, category Category) Category {
	queryInsert := "INSERT INTO categories(name) VALUES(?)"

	result, err := tx.ExecContext(ctx, queryInsert, category.Name)

	helpers.PanicError(err, "failed to exec query insert category")

	id, err := result.LastInsertId()

	helpers.PanicError(err, "failed to get last insert id category")

	createdCategory := repository.FindById(ctx, tx, int(id))

	return createdCategory
}

func (repository *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, limit, offset int) []Category {
	query := "SELECT id, name, created_at, updated_at FROM categories LIMIT ? OFFSET ?"

	rows, err := tx.QueryContext(ctx, query, limit, offset)

	helpers.PanicError(err, "failed to query all categories")

	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var category Category
		err := rows.Scan(&category.Id, &category.Name, &category.Created_At, &category.Updated_At)
		helpers.PanicError(err, "failed to scan all categories")

		categories = append(categories, category)
	}

	return categories
}

func (repository *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId int) Category {
	query := "SELECT id, name, created_at, updated_at FROM categories WHERE id = ?"

	rows, err := tx.QueryContext(ctx, query, categoryId)

	helpers.PanicError(err, "failed to query category by id")

	defer rows.Close()

	var category Category

	if rows.Next() {
		err := rows.Scan(&category.Id, &category.Name, &category.Created_At, &category.Updated_At)

		helpers.PanicError(err, "failed to scan category by id")
	} else {
		panic(exception.NewNotFoundError("category not found"))
	}

	return category
}

func (repository *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category Category) Category {
	query := "UPDATE categories SET name = ? WHERE id = ?"

	_, err := tx.ExecContext(ctx, query, category.Name, category.Id)

	helpers.PanicError(err, "failed to exec query update category")

	updatedCategory := repository.FindById(ctx, tx, category.Id)

	return updatedCategory
}

func (repository *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, categoryId int) {
	query := "DELETE FROM categories WHERE id = ?"

	result, err := tx.ExecContext(ctx, query, categoryId)

	helpers.PanicError(err, "failed to exec query delete category")

	resultRows, err := result.RowsAffected()

	if resultRows == 0 {
		panic(exception.NewNotFoundError("category not found"))
	}

	helpers.PanicError(err, "failed to display rows affected delete category")
}

func (repository *CategoryRepositoryImpl) CountCategories(ctx context.Context, tx *sql.Tx) int {
	query := "SELECT COUNT(*) FROM categories"

	rows, err := tx.QueryContext(ctx, query)

	helpers.PanicError(err, "failed to query count categories")

	defer rows.Close()

	var countCategory int

	if rows.Next() {
		err := rows.Scan(&countCategory)
		helpers.PanicError(err, "failed to scan count categories")
	}

	return countCategory
}
