package follow

import (
	"context"
	"database/sql"

	"github.com/hutamatr/GoBlogify/helpers"
)

type FollowService interface {
	Following(ctx context.Context, followerId, userId int) FollowResponse
	Unfollow(ctx context.Context, followerId, userId int)
	FindAllFollowed(ctx context.Context, userId, limit, offset int) []FollowJoinResponse
	FindAllFollower(ctx context.Context, userId, limit, offset int) []FollowJoinResponse
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

func (service *FollowServiceImpl) Following(ctx context.Context, followerId, userId int) FollowResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	newFollow := Follow{
		Follower_Id: followerId,
		Followed_Id: userId,
	}

	followingUser := service.repository.Save(ctx, tx, newFollow)

	return ToFollowResponse(followingUser)
}

func (services *FollowServiceImpl) Unfollow(ctx context.Context, followerId, userId int) {
	tx, err := services.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	services.repository.Delete(ctx, tx, followerId, userId)
}

func (service *FollowServiceImpl) FindAllFollowed(ctx context.Context, userId, limit, offset int) []FollowJoinResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	followed := service.repository.FindAllFollowedByUser(ctx, tx, userId, limit, offset)

	var followedData []FollowJoinResponse

	if len(followed) == 0 {
		return followedData
	}

	for _, follow := range followed {
		followedData = append(followedData, ToFollowJoinResponse(follow))
	}

	return followedData
}

func (service *FollowServiceImpl) FindAllFollower(ctx context.Context, userId, limit, offset int) []FollowJoinResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err, "failed to begin transaction")
	defer helpers.TxRollbackCommit(tx)

	followers := service.repository.FindAllFollowerByUser(ctx, tx, userId, limit, offset)

	var followerData []FollowJoinResponse

	if len(followers) == 0 {
		return followerData
	}

	for _, follower := range followers {
		followerData = append(followerData, ToFollowJoinResponse(follower))
	}

	return followerData
}
