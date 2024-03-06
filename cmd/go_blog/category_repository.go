package go_blog

import (
	"context"
	"database/sql"
	"errors"

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

func (a *CategoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, category Category) Category {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "INSERT INTO article(name) VALUES(?)"

	result, err := tx.ExecContext(ctxC, query, category.Name)

	helper.PanicError(err)

	id, err := result.LastInsertId()

	helper.PanicError(err)

	category.Id = int(id)

	return category
}

func (a *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []Category {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT id, name FROM category"

	rows, err := tx.QueryContext(ctxC, query)

	helper.PanicError(err)

	defer rows.Close()

	categories := []Category{}

	for rows.Next() {
		category := Category{}

		err := rows.Scan(&category.Id, &category.Name)

		helper.PanicError(err)

		categories = append(categories, category)
	}

	return categories
}

func (a *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId int) (Category, error) {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT id, name FROM category WHERE id = ? LIMIT = 1"

	rows, err := tx.QueryContext(ctxC, query, categoryId)

	helper.PanicError(err)

	defer rows.Close()

	category := Category{}

	if rows.Next() {
		err := rows.Scan(&category.Id, &category.Name)

		helper.PanicError(err)

		return category, nil
	} else {
		return category, errors.New("category not found")
	}
}

func (a *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category Category) Category {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "UPDATE category SET name = ? WHERE id = ?"

	result, err := tx.ExecContext(ctxC, query, category.Name, category.Id)

	helper.PanicError(err)

	resultRows, err := result.RowsAffected()

	helper.PanicError(err)

	if resultRows == 0 {
		panic(errors.New("category not found"))
	}

	return category
}

func (a *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, categoryId int) {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "DELETE FROM category WHERE id = ?"

	result, err := tx.ExecContext(ctxC, query, categoryId)

	helper.PanicError(err)

	resultRows, err := result.RowsAffected()

	helper.PanicError(err)

	if resultRows == 0 {
		panic(errors.New("category not found"))
	}
}
