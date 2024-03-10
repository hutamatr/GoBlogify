package go_blog

import (
	"context"
	"database/sql"

	"github.com/hutamatr/go-blog-api/cmd/go_blog/helper"
)

type CategoryRepository interface {
	Save(ctx context.Context, tx *sql.Tx, category Category) Category
	FindAll(ctx context.Context, tx *sql.Tx) []Category
	FindById(ctx context.Context, tx *sql.Tx, categoryId int) (Category, error)
	Update(ctx context.Context, tx *sql.Tx, category Category) Category
	Delete(ctx context.Context, tx *sql.Tx, categoryId int)
}

type CategoryRepositoryImpl struct {
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
}

func (repository *CategoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, category Category) Category {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	queryInsert := "INSERT INTO category(name) VALUES(?)"

	result, err := tx.ExecContext(ctxC, queryInsert, category.Name)

	helper.PanicError(err)

	id, err := result.LastInsertId()

	helper.PanicError(err)

	querySelect := "SELECT id, name, created_at, updated_at FROM category WHERE id = ?"

	rows, err := tx.QueryContext(ctxC, querySelect, id)

	helper.PanicError(err)

	defer rows.Close()

	createdCategory := Category{}

	if rows.Next() {
		err := rows.Scan(&createdCategory.Id, &createdCategory.Name, &createdCategory.Created_At, &createdCategory.Updated_At)
		helper.PanicError(err)
	}

	return createdCategory
}

func (repository *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []Category {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT id, name, created_at, updated_at FROM category"

	rows, err := tx.QueryContext(ctxC, query)

	helper.PanicError(err)

	defer rows.Close()

	categories := []Category{}

	for rows.Next() {
		category := Category{}

		err := rows.Scan(&category.Id, &category.Name, &category.Created_At, &category.Updated_At)

		helper.PanicError(err)

		categories = append(categories, category)
	}

	return categories
}

func (repository *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId int) (Category, error) {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT id, name, created_at, updated_at FROM category WHERE id = ?"

	rows, err := tx.QueryContext(ctxC, query, categoryId)

	helper.PanicError(err)

	defer rows.Close()

	category := Category{}

	if rows.Next() {
		err := rows.Scan(&category.Id, &category.Name, &category.Created_At, &category.Updated_At)

		helper.PanicError(err)

		return category, nil
	} else {
		return category, helper.NotFoundError
	}
}

func (repository *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category Category) Category {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "UPDATE category SET name = ? WHERE id = ?"

	result, err := tx.ExecContext(ctxC, query, category.Name, category.Id)

	helper.PanicError(err)

	resultRows, err := result.RowsAffected()

	helper.PanicError(err)

	if resultRows == 0 {
		panic(helper.NotFoundError)
	}

	return category
}

func (repository *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, categoryId int) {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "DELETE FROM category WHERE id = ?"

	result, err := tx.ExecContext(ctxC, query, categoryId)

	helper.PanicError(err)

	resultRows, err := result.RowsAffected()

	helper.PanicError(err)

	if resultRows == 0 {
		panic(helper.NotFoundError)
	}
}
