package role

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
)

type RoleService interface {
	Create(ctx context.Context, request RoleCreateRequest, isAdmin bool) RoleResponse
	FindAll(ctx context.Context, isAdmin bool) []RoleResponse
	FindById(ctx context.Context, roleId int, isAdmin bool) RoleResponse
	Update(ctx context.Context, request RoleUpdateRequest, isAdmin bool) RoleResponse
	Delete(ctx context.Context, roleId int, isAdmin bool)
}

type RoleServiceImpl struct {
	repository RoleRepository
	db         *sql.DB
	validator  *validator.Validate
}

func NewRoleService(roleRepository RoleRepository, db *sql.DB, validator *validator.Validate) RoleService {
	return &RoleServiceImpl{
		repository: roleRepository,
		db:         db,
		validator:  validator,
	}
}

func (service *RoleServiceImpl) Create(ctx context.Context, request RoleCreateRequest, isAdmin bool) RoleResponse {
	if !isAdmin {
		panic(exception.NewBadRequestError("only admin can create role"))
	}

	err := service.validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	roleRequest := Role{
		Name: request.Name,
	}

	createdRole := service.repository.Save(ctx, tx, roleRequest)

	return ToRoleResponse(createdRole)
}

func (service *RoleServiceImpl) FindAll(ctx context.Context, isAdmin bool) []RoleResponse {
	if !isAdmin {
		panic(exception.NewBadRequestError("only admin can get all roles"))
	}

	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	roles := service.repository.FindAll(ctx, tx)

	var rolesData []RoleResponse

	if len(roles) == 0 {
		panic(exception.NewNotFoundError("roles not found"))
	}

	for _, role := range roles {
		rolesData = append(rolesData, ToRoleResponse(role))
	}

	return rolesData
}

func (service *RoleServiceImpl) FindById(ctx context.Context, roleId int, isAdmin bool) RoleResponse {
	if !isAdmin {
		panic(exception.NewBadRequestError("only admin can get role by id"))
	}
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	role := service.repository.FindById(ctx, tx, roleId)

	return ToRoleResponse(role)
}

func (service *RoleServiceImpl) Update(ctx context.Context, request RoleUpdateRequest, isAdmin bool) RoleResponse {
	if !isAdmin {
		panic(exception.NewBadRequestError("only admin can update role"))
	}

	err := service.validator.Struct(request)
	helpers.PanicError(err, "invalid request")

	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	role := service.repository.FindById(ctx, tx, request.Id)

	role.Name = request.Name

	updatedRole := service.repository.Update(ctx, tx, role)

	return ToRoleResponse(updatedRole)
}

func (service *RoleServiceImpl) Delete(ctx context.Context, roleId int, isAdmin bool) {
	if !isAdmin {
		panic(exception.NewBadRequestError("only admin can delete role"))
	}
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	service.repository.Delete(ctx, tx, roleId)
}
