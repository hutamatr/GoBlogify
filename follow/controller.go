package follow

import (
	"net/http"
	"strconv"

	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/julienschmidt/httprouter"
)

type FollowController interface {
	FollowUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UnfollowUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllFollowedByUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllFollowerByUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type FollowControllerImpl struct {
	service FollowService
}

func NewFollowController(service FollowService) FollowController {
	return &FollowControllerImpl{
		service: service,
	}
}

func (controller *FollowControllerImpl) FollowUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("userId")
	userId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid User Id")

	id = params.ByName("toUserId")
	toUserId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid to User Id")

	controller.service.Following(request.Context(), userId, toUserId)

	followResponse := helpers.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
	}

	writer.WriteHeader(http.StatusCreated)
	helpers.EncodeJSONFromResponse(writer, followResponse)
}

func (controller *FollowControllerImpl) UnfollowUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("userId")
	userId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid User Id")

	id = params.ByName("toUserId")
	toUserId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid to User Id")

	controller.service.Unfollow(request.Context(), userId, toUserId)

	unfollowResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, unfollowResponse)
}

func (controller *FollowControllerImpl) FindAllFollowedByUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("userId")
	userId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid User Id")

	limit, offset := helpers.GetLimitOffset(request)

	followedData, countFollowed := controller.service.FindAllFollowed(request.Context(), userId, limit, offset)

	followedResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"followed": followedData,
			"limit":    limit,
			"offset":   offset,
			"total":    countFollowed,
		},
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, followedResponse)
}

func (controller *FollowControllerImpl) FindAllFollowerByUserHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("userId")
	userId, err := strconv.Atoi(id)
	helpers.PanicError(err, "Invalid User Id")

	limit, offset := helpers.GetLimitOffset(request)

	followerData, countFollower := controller.service.FindAllFollower(request.Context(), userId, limit, offset)

	followerResponse := helpers.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"follower": followerData,
			"limit":    limit,
			"offset":   offset,
			"total":    countFollower,
		},
	}

	writer.WriteHeader(http.StatusOK)
	helpers.EncodeJSONFromResponse(writer, followerResponse)
}
