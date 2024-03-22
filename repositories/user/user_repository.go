package users

import (
	"context"
	"database/sql"

	"github.com/hutamatr/GoBlogify/model/domain"
)

type UserRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user domain.User) domain.User
	FindAll(ctx context.Context, tx *sql.Tx) []domain.User
	FindOne(ctx context.Context, tx *sql.Tx, userId int, email string) domain.User
	Update(ctx context.Context, tx *sql.Tx, user domain.User) domain.User
	Delete(ctx context.Context, tx *sql.Tx, userId int)
	FindPassword(ctx context.Context, tx *sql.Tx, email string) string
}
