package role

import (
	"context"
	"database/sql"

	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
)

type RoleRepository interface {
	Save(ctx context.Context, tx *sql.Tx, role Role) Role
	FindAll(ctx context.Context, tx *sql.Tx) []Role
	FindById(ctx context.Context, tx *sql.Tx, roleId int) Role
	FindByName(ctx context.Context, tx *sql.Tx, roleName string) Role
	Update(ctx context.Context, tx *sql.Tx, role Role) Role
	Delete(ctx context.Context, tx *sql.Tx, roleId int)
}

type RoleRepositoryImpl struct {
}

func NewRoleRepository() RoleRepository {
	return &RoleRepositoryImpl{}
}

func (repository *RoleRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, role Role) Role {
	query := "INSERT INTO role(name) VALUES (?)"

	result, err := tx.ExecContext(ctx, query, role.Name)

	helpers.PanicError(err, "failed to exec query insert role")

	id, err := result.LastInsertId()

	helpers.PanicError(err, "failed to get last insert id role")

	newRole := repository.FindById(ctx, tx, int(id))

	return newRole
}

func (repository *RoleRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []Role {
	query := "SELECT id, name, created_at, updated_at FROM role"

	rows, err := tx.QueryContext(ctx, query)

	helpers.PanicError(err, "failed to query all roles")

	defer rows.Close()

	var roles []Role

	for rows.Next() {
		var role Role
		err := rows.Scan(&role.Id, &role.Name, &role.Created_At, &role.Updated_At)
		helpers.PanicError(err, "failed to scan all roles")

		roles = append(roles, role)
	}

	return roles
}

func (repository *RoleRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, roleId int) Role {
	query := "SELECT id, name, created_at, updated_at FROM role WHERE id = ?"

	rows, err := tx.QueryContext(ctx, query, roleId)

	helpers.PanicError(err, "failed to query role by id")

	defer rows.Close()

	var role Role

	if rows.Next() {
		err := rows.Scan(&role.Id, &role.Name, &role.Created_At, &role.Updated_At)
		helpers.PanicError(err, "failed to scan role by id")
	} else {
		panic(exception.NewNotFoundError("role not found"))
	}

	return role
}

func (repository *RoleRepositoryImpl) FindByName(ctx context.Context, tx *sql.Tx, roleName string) Role {
	query := "SELECT id, name, created_at, updated_at FROM role WHERE name = ?"

	rows, err := tx.QueryContext(ctx, query, roleName)

	helpers.PanicError(err, "failed to query role by name")

	defer rows.Close()

	var role Role

	if rows.Next() {
		err := rows.Scan(&role.Id, &role.Name, &role.Created_At, &role.Updated_At)
		helpers.PanicError(err, "failed to scan role by name")
	}
	return role
}

func (repository *RoleRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, role Role) Role {
	query := "UPDATE role SET name = ? WHERE id = ?"

	_, err := tx.ExecContext(ctx, query, role.Name, role.Id)

	helpers.PanicError(err, "failed to exec query update role")

	updatedRole := repository.FindById(ctx, tx, role.Id)

	return updatedRole
}

func (repository *RoleRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, roleId int) {
	query := "DELETE FROM role WHERE id = ?"

	result, err := tx.ExecContext(ctx, query, roleId)

	helpers.PanicError(err, "failed to exec query delete role")

	resultRows, err := result.RowsAffected()

	if resultRows == 0 {
		panic(exception.NewNotFoundError("role not found"))
	}

	helpers.PanicError(err, "failed to display rows affected delete role")
}
