package follow

import (
	"context"
	"database/sql"

	"github.com/hutamatr/GoBlogify/exception"
	"github.com/hutamatr/GoBlogify/helpers"
)

type FollowRepositories interface {
	Save(ctx context.Context, tx *sql.Tx, follow Follow) Follow
	FindAllFollowerByUser(ctx context.Context, tx *sql.Tx, followedId, limit, offset int) []FollowJoin
	FindAllFollowedByUser(ctx context.Context, tx *sql.Tx, followerId, limit, offset int) []FollowJoin
	FindById(ctx context.Context, tx *sql.Tx, followId int) Follow
	Delete(ctx context.Context, tx *sql.Tx, followerId, followedId int)
}

type FollowRepositoriesImpl struct {
}

func NewFollowRepository() *FollowRepositoriesImpl {
	return &FollowRepositoriesImpl{}
}

func (repository *FollowRepositoriesImpl) Save(ctx context.Context, tx *sql.Tx, follow Follow) Follow {
	queryInsert := "INSERT INTO follow(follower_id, followed_id) VALUES(?, ?)"

	result, err := tx.ExecContext(ctx, queryInsert, follow.Follower_Id, follow.Followed_Id)

	helpers.PanicError(err, "failed to exec query insert follow")

	id, err := result.LastInsertId()

	helpers.PanicError(err, "failed to get last insert id follow")

	followData := repository.FindById(ctx, tx, int(id))

	return followData
}

func (repository *FollowRepositoriesImpl) FindAllFollowerByUser(ctx context.Context, tx *sql.Tx, followedId, limit, offset int) []FollowJoin {
	query := "SELECT user.id, user.username, user.first_name, user.last_name, follow.id, follow.followed_id, follow.follower_id, follow.created_at, follow.updated_at FROM follow JOIN user ON follow.follower_id = user.id WHERE follow.followed_id = ? LIMIT ? OFFSET ?"

	rows, err := tx.QueryContext(ctx, query, followedId, limit, offset)

	helpers.PanicError(err, "failed to query all follower by user")

	defer rows.Close()

	var followers []FollowJoin
	var firstName sql.NullString
	var lastName sql.NullString

	if rows.Next() {
		var follower FollowJoin
		err := rows.Scan(&follower.User.Id, &follower.User.Username, &firstName, &lastName, &follower.Id, &follower.Followed_Id, &follower.Follower_Id, &follower.Created_At, &follower.Updated_At)
		helpers.PanicError(err, "failed to scan all follower by user")

		if firstName.Valid {
			follower.User.First_Name = firstName.String
		} else {
			follower.User.First_Name = ""
		}

		if lastName.Valid {
			follower.User.Last_Name = lastName.String
		} else {
			follower.User.Last_Name = ""
		}

		followers = append(followers, follower)
	}

	return followers
}

func (repository *FollowRepositoriesImpl) FindAllFollowedByUser(ctx context.Context, tx *sql.Tx, followerId, limit, offset int) []FollowJoin {
	query := "SELECT user.id, user.username, user.first_name, user.last_name, follow.id, follow.followed_id, follow.follower_id, follow.created_at, follow.updated_at FROM follow JOIN user ON follow.followed_id = user.id WHERE follow.follower_id = ? LIMIT ? OFFSET ?"

	rows, err := tx.QueryContext(ctx, query, followerId, limit, offset)

	helpers.PanicError(err, "failed to query all followed by user")

	defer rows.Close()

	var followed []FollowJoin
	var firstName sql.NullString
	var lastName sql.NullString

	if rows.Next() {
		var follow FollowJoin
		err := rows.Scan(&follow.User.Id, &follow.User.Username, &firstName, &lastName, &follow.Id, &follow.Followed_Id, &follow.Follower_Id, &follow.Created_At, &follow.Updated_At)
		helpers.PanicError(err, "failed to scan all followed by user")

		if firstName.Valid {
			follow.User.First_Name = firstName.String
		} else {
			follow.User.First_Name = ""
		}

		if lastName.Valid {
			follow.User.Last_Name = lastName.String
		} else {
			follow.User.Last_Name = ""
		}

		followed = append(followed, follow)
	}

	return followed
}

func (repository *FollowRepositoriesImpl) FindById(ctx context.Context, tx *sql.Tx, followId int) Follow {
	query := "SELECT id, follower_id, followed_id FROM follow WHERE id = ?"

	rows, err := tx.QueryContext(ctx, query, followId)

	helpers.PanicError(err, "failed to query follow by id")

	defer rows.Close()

	var follow Follow

	if rows.Next() {
		err := rows.Scan(&follow.Id, &follow.Follower_Id, &follow.Followed_Id)
		helpers.PanicError(err, "failed to scan follow by id")
	}

	return follow
}

func (repository *FollowRepositoriesImpl) Delete(ctx context.Context, tx *sql.Tx, followerId, userId int) {
	query := "DELETE FROM follow WHERE follower_id = ? AND followed_id = ?"

	result, err := tx.ExecContext(ctx, query, followerId, userId)
	helpers.PanicError(err, "failed to exec query delete follow")

	resultRows, err := result.RowsAffected()
	helpers.PanicError(err, "failed to display rows affected delete follow")

	if resultRows == 0 {
		panic(exception.NewNotFoundError("follow not found"))
	}
}
