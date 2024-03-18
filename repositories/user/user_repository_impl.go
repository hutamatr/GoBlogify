package users

import (
	"context"
	"database/sql"
	"time"

	"github.com/hutamatr/go-blog-api/exception"
	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/domain"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	queryInsert := "INSERT INTO user(username, email, password, role_id) VALUES (?, ?, ?)"

	result, err := tx.ExecContext(ctxC, queryInsert, user.Username, user.Email, user.Password, user.Role_Id)

	helpers.PanicError(err)

	id, err := result.LastInsertId()

	helpers.PanicError(err)

	newUser := repository.FindOne(ctx, tx, int(id), user.Email)

	return newUser
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.User {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT id, username, email,first_name, last_name, role, created_at, updated_at, deleted_at FROM user LIMIT 10"

	rows, err := tx.QueryContext(ctxC, query)

	helpers.PanicError(err)

	defer rows.Close()

	var users []domain.User

	var deletedAt sql.NullTime

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.First_Name, &user.Last_Name, &user.Role_Id, &user.Created_At, &user.Updated_At, &deletedAt)

		helpers.PanicError(err)

		if deletedAt.Valid {
			user.Deleted_At = deletedAt.Time
		} else {
			user.Deleted_At = time.Time{}
		}

		users = append(users, user)
	}

	return users
}

func (repository *UserRepositoryImpl) FindOne(ctx context.Context, tx *sql.Tx, userId int, email string) domain.User {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	var rows *sql.Rows
	var err error

	if userId > 0 {
		query := "SELECT id, username, email, first_name, last_name, role, created_at, updated_at, deleted_at FROM user WHERE id = ?"
		rows, err = tx.QueryContext(ctxC, query, userId)
		helpers.PanicError(err)
	} else if email != "" {
		query := "SELECT id, username, email, first_name, last_name, role, created_at, updated_at, deleted_at FROM user WHERE email = ?"
		rows, err = tx.QueryContext(ctxC, query, email)
		helpers.PanicError(err)
	}

	defer rows.Close()

	var user domain.User

	var deletedAt sql.NullTime

	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.First_Name, &user.Last_Name, &user.Role_Id, &user.Created_At, &user.Updated_At, &deletedAt)

		helpers.PanicError(err)

		if deletedAt.Valid {
			user.Deleted_At = deletedAt.Time
		} else {
			user.Deleted_At = time.Time{}
		}
	} else {
		panic(exception.NewNotFoundError("user not found"))
	}

	return user
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	var result sql.Result
	var err error

	if user.Id > 0 {
		query := "UPDATE user SET username = ?, email = ?, first_name = ?, last_name = ? WHERE id = ?"
		result, err = tx.ExecContext(ctxC, query, user.Username, user.Email, user.First_Name, user.Last_Name, user.Id)
		helpers.PanicError(err)
	} else if user.Email != "" {
		query := "UPDATE user SET username = ?, email = ?, first_name = ?, last_name = ? WHERE email = ?"
		result, err = tx.ExecContext(ctxC, query, user.Username, user.Email, user.First_Name, user.Last_Name, user.Email)
		helpers.PanicError(err)
	}

	id, err := result.LastInsertId()

	helpers.PanicError(err)

	updatedUser := repository.FindOne(ctx, tx, int(id), user.Email)

	return updatedUser
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, userId int) {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "DELETE FROM user WHERE id = ?"
	_, err := tx.ExecContext(ctxC, query, userId)
	helpers.PanicError(err)
}

func (repository *UserRepositoryImpl) FindPassword(ctx context.Context, tx *sql.Tx, email string) string {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT password FROM user WHERE email = ?"
	rows, err := tx.QueryContext(ctxC, query, email)
	helpers.PanicError(err)

	defer rows.Close()

	var password string

	if rows.Next() {
		err := rows.Scan(&password)
		helpers.PanicError(err)
	} else {
		panic(exception.NewNotFoundError("user not found"))
	}

	return password
}
