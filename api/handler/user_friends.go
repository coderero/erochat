package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/coderero/erochat-server/interfaces"
	"github.com/coderero/erochat-server/types"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserFriendShipHandler struct {
	// userStore is a data store for user.
	userStore interfaces.UserStore

	// friendStore is a data store for friend.
	friendStore interfaces.FriendStore

	// validate is a validator that validates the request.
	validate *validator.Validate
}

// NewUserFriendShipHandler returns a new user friend ship handler.
func NewUserFriendShipHandler(validator *validator.Validate, userStore interfaces.UserStore, friendStore interfaces.FriendStore) *UserFriendShipHandler {
	return &UserFriendShipHandler{
		userStore:   userStore,
		friendStore: friendStore,
		validate:    validator,
	}
}

var (
	sww = &echo.HTTPError{
		Code:    echo.ErrBadRequest.Code,
		Message: "something went wrong",
	}
	fnf = &echo.HTTPError{
		Code:    echo.ErrNotFound.Code,
		Message: "friend not found",
	}
)

// GetFriends gets the friends of a user.
func (u *UserFriendShipHandler) GetFriends(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok {
		return sww
	}

	// Parse the user id.
	userID, err := uuid.Parse(uid)
	if err != nil {
		return sww
	}

	friends, err := u.friendStore.GetFriends(userID)
	if err != nil {
		return sww
	}

	res := types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "friends fetched successfully",
		Data:    friends,
	}

	return c.JSON(http.StatusOK, res)
}

// GetFriend gets a friend by its id.
func (u *UserFriendShipHandler) GetFriend(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok {
		return sww
	}

	// Parse the user id.
	userID, err := uuid.Parse(uid)
	if err != nil {
		return sww
	}

	fid := c.Param("uid")
	if len(fid) == 0 {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "uid is required in url parameter",
		}
	}

	friendID, err := uuid.Parse(fid)
	if err != nil {
		fmt.Println(err)
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "invalid user id",
		}
	}

	friend, err := u.friendStore.GetFriend(userID, friendID)
	if err != nil {
		if errors.Is(err, interfaces.ErrFriendNotFound) {
			return fnf
		}
		return sww
	}

	res := types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "friend fetched successfully",
		Data:    friend,
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteFriend deletes a friend by its id.
func (u *UserFriendShipHandler) DeleteFriend(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok {
		return sww
	}

	// Parse the user id.
	userID, err := uuid.Parse(uid)
	if err != nil {
		return sww
	}

	fid := c.Param("uid")
	if len(fid) == 0 {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "uid is required in url parameter",
		}
	}
	friendID, err := uuid.Parse(fid)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "invalid user id",
		}
	}

	err = u.friendStore.DeleteFriend(userID, friendID)
	if err != nil {
		if errors.Is(err, interfaces.ErrFriendNotFound) {
			return fnf
		}
		return sww
	}

	res := types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "friend deleted successfully",
	}

	return c.JSON(http.StatusOK, res)
}

// GetFriendRequests gets the friend requests of a user.
func (u *UserFriendShipHandler) GetFriendRequests(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok {
		return sww
	}

	// Parse the user id.
	userID, err := uuid.Parse(uid)
	if err != nil {
		return sww
	}

	requests, err := u.friendStore.GetFriendRequests(userID)
	if err != nil {
		return sww
	}

	res := types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "friend requests fetched successfully",
		Data:    requests,
	}

	return c.JSON(http.StatusOK, res)
}

// GetFriendRequest gets a friend request by its id.
func (u *UserFriendShipHandler) GetFriendRequest(c echo.Context) error {
	uUID, ok := c.Get("uid").(string)
	if !ok {
		return sww
	}

	// Parse the user id.
	userID, err := uuid.Parse(uUID)
	if err != nil {
		return sww
	}

	fid := c.Param("uid")
	if len(fid) == 0 {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "uid is required in url parameter",
		}
	}
	reqId, err := uuid.Parse(fid)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "invalid user id",
		}
	}

	request, err := u.friendStore.GetFriendRequest(userID, reqId)
	if err != nil {
		if errors.Is(err, interfaces.ErrFriendNotFound) {
			return fnf
		}
		return sww
	}

	res := types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "friend request fetched successfully",
		Data:    request,
	}

	return c.JSON(http.StatusOK, res)
}

// AcceptFriendRequest accepts a friend request.
func (u *UserFriendShipHandler) AcceptFriendRequest(c echo.Context) error {
	rUUID, ok := c.Get("uid").(string)
	if !ok {
		return sww
	}

	UUID, err := uuid.Parse(rUUID)
	if err != nil {
		return sww
	}

	fid := c.Param("uid")
	if len(fid) == 0 {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "uid is required in url parameter",
		}
	}
	reqId, err := uuid.Parse(fid)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "invalid user id",
		}
	}

	err = u.friendStore.AcceptFriendRequest(UUID, reqId)
	if err != nil {
		if errors.Is(err, interfaces.ErrFriendNotFound) {
			return fnf
		}
		return sww
	}

	res := types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "friend request accepted successfully",
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteFriendRequest deletes a friend request by its id.
func (u *UserFriendShipHandler) DeleteFriendRequest(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok {
		return sww
	}

	// Parse the user id.
	userID, err := uuid.Parse(uid)
	if err != nil {
		return sww
	}

	fid := c.Param("uid")
	if len(fid) == 0 {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "uid is required in url parameter",
		}
	}
	relationID, err := uuid.Parse(fid)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "invalid user id",
		}
	}

	err = u.friendStore.DeleteFriendRequest(userID, relationID)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, interfaces.ErrFriendNotFound) {
			return fnf
		}
		return sww
	}

	res := types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "friend request deleted successfully",
	}

	return c.JSON(http.StatusOK, res)
}

// GetFriendsStatus gets the friends status of a user.
func (u *UserFriendShipHandler) GetFriendsStatus(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok {
		return sww
	}

	// Parse the user id.
	userID, err := uuid.Parse(uid)
	if err != nil {
		return sww
	}

	statu, err := u.friendStore.GetFriendsStatus(userID)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return sww
	}

	res := types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "friends status fetched successfully",
		Data:    statu,
	}

	return c.JSON(http.StatusOK, res)
}

// GetFriendStatus gets the friend status of a user.
func (u *UserFriendShipHandler) GetFriendStatus(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok {
		return sww
	}

	// Parse the user id.
	userID, err := uuid.Parse(uid)
	if err != nil {
		return sww
	}

	fid := c.Param("uid")
	if len(fid) == 0 {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "uid is required in url parameter",
		}
	}
	friendID, err := uuid.Parse(fid)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "invalid user id",
		}
	}

	status, err := u.friendStore.GetFriendStatus(userID, friendID)
	if err != nil {
		if errors.Is(err, interfaces.ErrFriendStatusNotFound) {
			return &echo.HTTPError{
				Code:    echo.ErrNotFound.Code,
				Message: "friend status not found",
			}
		}
		return sww
	}

	res := types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "friend status fetched successfully",
		Data:    status,
	}

	return c.JSON(http.StatusOK, res)
}
