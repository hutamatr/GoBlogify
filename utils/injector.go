//go:build wireinject
// +build wireinject

package utils

import (
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/hutamatr/GoBlogify/admin"
	"github.com/hutamatr/GoBlogify/category"
	"github.com/hutamatr/GoBlogify/comment"
	"github.com/hutamatr/GoBlogify/follow"
	"github.com/hutamatr/GoBlogify/post"
	"github.com/hutamatr/GoBlogify/role"
	"github.com/hutamatr/GoBlogify/user"
)

func InitializedRoleController(db *sql.DB, validator *validator.Validate) role.RoleController {
	wire.Build(role.NewRoleRepository, role.NewRoleService, role.NewRoleController)
	return nil
}

func InitializedUserController(db *sql.DB, validator *validator.Validate) user.UserController {
	wire.Build(user.NewUserRepository, user.NewUserService, user.NewUserController, role.NewRoleRepository)
	return nil
}

func InitializedAdminController(db *sql.DB, validator *validator.Validate) admin.AdminController {
	wire.Build(admin.NewAdminService, admin.NewAdminController, role.NewRoleRepository, user.NewUserRepository)
	return nil
}

func InitializedPostController(db *sql.DB, validator *validator.Validate) post.PostController {
	wire.Build(post.NewPostRepository, post.NewPostService, post.NewPostController)
	return nil
}

func InitializedCommentController(db *sql.DB, validator *validator.Validate) comment.CommentController {
	wire.Build(comment.NewCommentRepository, comment.NewCommentService, comment.NewCommentController)
	return nil
}

func InitializedCategoryController(db *sql.DB, validator *validator.Validate) category.CategoryController {
	wire.Build(category.NewCategoryRepository, category.NewCategoryService, category.NewCategoryController)
	return nil
}

func InitializedFollowController(db *sql.DB) follow.FollowController {
	wire.Build(follow.NewFollowRepository, follow.NewFollowService, follow.NewFollowController)
	return nil
}
