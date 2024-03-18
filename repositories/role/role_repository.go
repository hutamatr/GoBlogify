package role

import (
	"context"
	"database/sql"

	"github.com/hutamatr/go-blog-api/model/domain"
)

type RoleRepository interface {
	Save(ctx context.Context, tx *sql.Tx, role domain.Role) domain.Role
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Role
	FindById(ctx context.Context, tx *sql.Tx, roleId int) domain.Role
	Update(ctx context.Context, tx *sql.Tx, role domain.Role) domain.Role
	Delete(ctx context.Context, tx *sql.Tx, roleId int)
}
