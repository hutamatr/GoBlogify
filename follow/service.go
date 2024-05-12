package follow

import (
	"context"
	"database/sql"

	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
)

type FollowService interface {
	Following(ctx context.Context, userId, toUserId int) FollowResponse
	Unfollow(ctx context.Context, userId, toUserId int)
	FindAllFollowed(ctx context.Context, userId, limit, offset int) ([]FollowJoinResponse, int)
	FindAllFollower(ctx context.Context, userId, limit, offset int) ([]FollowJoinResponse, int)
}

type FollowServiceImpl struct {
	repository FollowRepositories
	db         *sql.DB
}

func NewFollowService(repository FollowRepositories, db *sql.DB) FollowService {
	return &FollowServiceImpl{
		repository: repository,
		db:         db,
	}
}

func (service *FollowServiceImpl) Following(ctx context.Context, userId, toUserId int) FollowResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	newFollow := Follow{
		Follower_Id: userId,
		Followed_Id: toUserId,
	}

	followingUser := service.repository.Save(ctx, tx, newFollow)

	return ToFollowResponse(followingUser)
}

func (services *FollowServiceImpl) Unfollow(ctx context.Context, userId, toUserId int) {
	tx, err := services.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	services.repository.Delete(ctx, tx, userId, toUserId)
}

func (service *FollowServiceImpl) FindAllFollowed(ctx context.Context, userId, limit, offset int) ([]FollowJoinResponse, int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	followed := service.repository.FindAllFollowedByUser(ctx, tx, userId, limit, offset)
	countFollowed := service.repository.CountFollowed(ctx, tx, userId)

	var followedData []FollowJoinResponse

	if len(followed) == 0 {
		panic(exception.NewNotFoundError("followed not found"))
	}

	for _, follow := range followed {
		followedData = append(followedData, ToFollowJoinResponse(follow))
	}

	return followedData, countFollowed
}

func (service *FollowServiceImpl) FindAllFollower(ctx context.Context, userId, limit, offset int) ([]FollowJoinResponse, int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	followers := service.repository.FindAllFollowerByUser(ctx, tx, userId, limit, offset)
	countFollower := service.repository.CountFollower(ctx, tx, userId)

	var followerData []FollowJoinResponse

	if len(followers) == 0 {
		panic(exception.NewNotFoundError("follower not found"))
	}

	for _, follower := range followers {
		followerData = append(followerData, ToFollowJoinResponse(follower))
	}

	return followerData, countFollower
}
