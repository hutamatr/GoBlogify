package role

import (
	"context"
	"database/sql"

	"github.com/hutamatr/go-blog-api/exception"
	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/domain"
)

type RoleRepositoryImpl struct {
}

func NewRoleRepository() RoleRepository {
	return &RoleRepositoryImpl{}
}

func (repository *RoleRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, role domain.Role) domain.Role {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "INSERT INTO role(name) VALUES (?)"

	result, err := tx.ExecContext(ctxC, query, role.Name)

	helpers.PanicError(err)

	id, err := result.LastInsertId()

	helpers.PanicError(err)

	newRole := repository.FindById(ctx, tx, int(id))

	return newRole
}

func (repository *RoleRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Role {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT id, name, created_at, updated_at FROM role"

	rows, err := tx.QueryContext(ctxC, query)

	helpers.PanicError(err)

	defer rows.Close()

	var roles []domain.Role

	for rows.Next() {
		var role domain.Role
		err := rows.Scan(&role.Id, &role.Name, &role.Created_At, &role.Updated_At)
		helpers.PanicError(err)

		roles = append(roles, role)
	}

	return roles
}

func (repository *RoleRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, roleId int) domain.Role {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT id, name, created_at, updated_at FROM role WHERE id = ?"

	rows, err := tx.QueryContext(ctxC, query, roleId)

	helpers.PanicError(err)

	defer rows.Close()

	var role domain.Role

	if rows.Next() {
		err := rows.Scan(&role.Id, &role.Name, &role.Created_At, &role.Updated_At)
		helpers.PanicError(err)
	} else {
		panic(exception.NewNotFoundError("role not found"))
	}

	return role
}

func (repository *RoleRepositoryImpl) FindByName(ctx context.Context, tx *sql.Tx, roleName string) domain.Role {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "SELECT id, name, created_at, updated_at FROM role WHERE name = ?"

	rows, err := tx.QueryContext(ctxC, query, roleName)

	helpers.PanicError(err)

	defer rows.Close()

	var role domain.Role

	if rows.Next() {
		err := rows.Scan(&role.Id, &role.Name, &role.Created_At, &role.Updated_At)
		helpers.PanicError(err)
	}
	return role
}

func (repository *RoleRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, role domain.Role) domain.Role {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "UPDATE role SET name = ? WHERE id = ?"

	_, err := tx.ExecContext(ctxC, query, role.Name, role.Id)

	helpers.PanicError(err)

	updatedRole := repository.FindById(ctx, tx, role.Id)

	return updatedRole
}

func (repository *RoleRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, roleId int) {
	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	query := "DELETE FROM role WHERE id = ?"

	result, err := tx.ExecContext(ctxC, query, roleId)

	helpers.PanicError(err)

	resultRows, err := result.RowsAffected()

	if resultRows == 0 {
		panic(exception.NewNotFoundError("role not found"))
	}

	helpers.PanicError(err)
}
