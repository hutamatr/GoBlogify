package repositories

import (
	"context"
	"database/sql"
	"log"

	"github.com/hutamatr/go-blog-api/exception"
	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/domain"
)

type CategoryRepositoryImpl struct {
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
}

func (repository *CategoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	queryInsert := "INSERT INTO category(name) VALUES(?)"

	result, err := tx.ExecContext(ctxC, queryInsert, category.Name)

	helpers.PanicError(err)

	id, err := result.LastInsertId()

	helpers.PanicError(err)

	querySelect := "SELECT id, name, created_at, updated_at FROM category WHERE id = ?"

	rows, err := tx.QueryContext(ctxC, querySelect, id)

	helpers.PanicError(err)

	defer rows.Close()

	var createdCategory domain.Category

	if rows.Next() {
		err := rows.Scan(&createdCategory.Id, &createdCategory.Name, &createdCategory.Created_At, &createdCategory.Updated_At)
		helpers.PanicError(err)
	}

	return createdCategory
}

func (repository *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Category {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT id, name, created_at, updated_at FROM category"

	rows, err := tx.QueryContext(ctxC, query)

	helpers.PanicError(err)

	defer rows.Close()

	var categories []domain.Category

	for rows.Next() {
		var category domain.Category

		err := rows.Scan(&category.Id, &category.Name, &category.Created_At, &category.Updated_At)

		helpers.PanicError(err)

		categories = append(categories, category)
	}

	return categories
}

func (repository *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId int) domain.Category {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT id, name, created_at, updated_at FROM category WHERE id = ?"

	rows, err := tx.QueryContext(ctxC, query, categoryId)

	helpers.PanicError(err)

	defer rows.Close()

	var category domain.Category

	if rows.Next() {
		err := rows.Scan(&category.Id, &category.Name, &category.Created_At, &category.Updated_At)

		helpers.PanicError(err)
	} else {
		panic(exception.NewNotFoundError("category not found"))
	}

	return category
}

func (repository *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "UPDATE category SET name = ? WHERE id = ?"

	result, err := tx.ExecContext(ctxC, query, category.Name, category.Id)

	helpers.PanicError(err)

	resultRows, err := result.RowsAffected()

	helpers.PanicError(err)

	if resultRows == 0 {
		log.Fatalf("expected single row affected, got %d rows affected", resultRows)
		panic(exception.NewNotFoundError("category not found"))
	}

	return category
}

func (repository *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, categoryId int) {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "DELETE FROM category WHERE id = ?"

	_, err := tx.ExecContext(ctxC, query, categoryId)

	helpers.PanicError(err)

}
