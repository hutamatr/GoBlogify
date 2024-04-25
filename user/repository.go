package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
)

type UserRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user User) UserJoin
	FindAll(ctx context.Context, tx *sql.Tx) []UserJoin
	FindOne(ctx context.Context, tx *sql.Tx, userId int, email string) UserJoin
	Update(ctx context.Context, tx *sql.Tx, user UserJoin) UserJoin
	Delete(ctx context.Context, tx *sql.Tx, userId int)
	FindPassword(ctx context.Context, tx *sql.Tx, email string) string
}

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user User) UserJoin {
	queryInsert := "INSERT INTO user(username, email, password, role_id) VALUES (?, ?, ?, ?)"

	result, err := tx.ExecContext(ctx, queryInsert, user.Username, user.Email, user.Password, user.Role_Id)

	helpers.PanicError(err, "failed to exec query insert user")

	id, err := result.LastInsertId()

	helpers.PanicError(err, "failed to get last insert id user")

	newUser := repository.FindOne(ctx, tx, int(id), "")

	return newUser
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []UserJoin {
	query := `SELECT u.id, u.username, u.email, u.first_name, u.last_name, u.role_id, u.created_at, u.updated_at, u.deleted_at,
	(SELECT COUNT(*) FROM follow f WHERE f.followed_id = u.id) AS follower_count,
	(SELECT COUNT(*) FROM follow f WHERE f.follower_id = u.id) AS following_count
	FROM user u WHERE u.is_deleted = false LIMIT 10`

	rows, err := tx.QueryContext(ctx, query)

	helpers.PanicError(err, "failed to query all users")

	defer rows.Close()

	var users []UserJoin

	var deletedAt sql.NullTime
	var firstName sql.NullString
	var lastName sql.NullString

	for rows.Next() {
		var user UserJoin
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &firstName, &lastName, &user.Role_Id, &user.Created_At, &user.Updated_At, &deletedAt, &user.Follower, &user.Following)

		helpers.PanicError(err, "failed to scan all users")

		if deletedAt.Valid {
			user.Deleted_At = deletedAt.Time
		} else {
			user.Deleted_At = time.Time{}
		}

		if firstName.Valid {
			user.First_Name = firstName.String
		} else {
			user.First_Name = ""
		}

		if lastName.Valid {
			user.Last_Name = lastName.String
		} else {
			user.Last_Name = ""
		}

		users = append(users, user)
	}

	return users
}

func (repository *UserRepositoryImpl) FindOne(ctx context.Context, tx *sql.Tx, userId int, email string) UserJoin {
	var rows *sql.Rows
	var err error

	if userId > 0 {
		query := `SELECT u.id, u.username, u.email, u.first_name, u.last_name, u.role_id, u.created_at, u.updated_at, u.deleted_at,
		(SELECT COUNT(*) FROM follow f WHERE f.followed_id = u.id) AS follower_count,
		(SELECT COUNT(*) FROM follow f WHERE f.follower_id = u.id) AS following_count
		FROM user u WHERE u.id = ? AND u.is_deleted = false`

		rows, err = tx.QueryContext(ctx, query, userId)
		helpers.PanicError(err, "failed to query one user")
	} else if email != "" {
		query := `SELECT u.id, u.username, u.email, u.first_name, u.last_name, u.role_id, u.created_at, u.updated_at, u.deleted_at,
		(SELECT COUNT(*) FROM follow f WHERE f.followed_id = u.id) AS follower_count,
		(SELECT COUNT(*) FROM follow f WHERE f.follower_id = u.id) AS following_count
		FROM user u WHERE u.email = ? AND u.is_deleted = false`

		rows, err = tx.QueryContext(ctx, query, email)
		helpers.PanicError(err, "failed to query one user")
	} else {
		panic(exception.NewNotFoundError("user not found"))
	}

	defer rows.Close()

	var user UserJoin

	var deletedAt sql.NullTime

	var firstName sql.NullString
	var lastName sql.NullString

	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &firstName, &lastName, &user.Role_Id, &user.Created_At, &user.Updated_At, &deletedAt, &user.Follower, &user.Following)

		helpers.PanicError(err, "failed to scan one user")

		if deletedAt.Valid {
			user.Deleted_At = deletedAt.Time
		} else {
			user.Deleted_At = time.Time{}
		}

		if firstName.Valid {
			user.First_Name = firstName.String
		} else {
			user.First_Name = ""
		}

		if lastName.Valid {
			user.Last_Name = lastName.String
		} else {
			user.Last_Name = ""
		}
	}

	return user
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user UserJoin) UserJoin {
	var result sql.Result
	var err error

	if user.Id > 0 {
		query := "UPDATE user SET username = ?, first_name = ?, last_name = ? WHERE id = ? AND is_deleted = false"
		result, err = tx.ExecContext(ctx, query, user.Username, user.First_Name, user.Last_Name, user.Id)
		helpers.PanicError(err, "failed to exec query update user")
	} else if user.Email != "" {
		query := "UPDATE user SET username = ?, first_name = ?, last_name = ? WHERE email = ? AND is_deleted = false"
		result, err = tx.ExecContext(ctx, query, user.Username, user.First_Name, user.Last_Name, user.Email)
		helpers.PanicError(err, "failed to exec query update user")
	}

	id, err := result.LastInsertId()

	helpers.PanicError(err, "failed to get last insert id user")

	updatedUser := repository.FindOne(ctx, tx, int(id), user.Email)

	return updatedUser
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, userId int) {
	query := "UPDATE user SET deleted_at = NOW(), is_deleted = true WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, userId)
	helpers.PanicError(err, "failed to exec query delete user")
}

func (repository *UserRepositoryImpl) FindPassword(ctx context.Context, tx *sql.Tx, email string) string {
	query := "SELECT password FROM user WHERE email = ? AND is_deleted = false"
	rows, err := tx.QueryContext(ctx, query, email)
	helpers.PanicError(err, "failed to query password user")

	defer rows.Close()

	var password string

	if rows.Next() {
		err := rows.Scan(&password)
		helpers.PanicError(err, "failed to scan password user")
	} else {
		panic(exception.NewNotFoundError("user not found"))
	}

	return password
}
