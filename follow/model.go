package follow

import (
	"time"

	"github.com/hutamatr/GoBlogify/user"
)

type Follow struct {
	Id          int
	Follower_Id int
	Followed_Id int
	Created_At  time.Time
	Updated_At  time.Time
}

type FollowJoin struct {
	Id          int
	Follower_Id int
	Followed_Id int
	Created_At  time.Time
	Updated_At  time.Time
	User        user.User
}
