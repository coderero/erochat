package handler

import (
	"github.com/coderero/erochat-server/interfaces"
	"github.com/go-playground/validator/v10"
)

type UserFirendShipHandler struct {
	// userStore is a data store for user.
	userStore interfaces.UserStore

	// friendStore is a data store for friend.
	friendStore interfaces.FriendStore

	// validate is a validator that validates the request.
	validate *validator.Validate
}

// NewUserFriendShipHandler returns a new user friend ship handler.
func NewUserFriendShipHandler(validator *validator.Validate, userStore interfaces.UserStore, friendStore interfaces.FriendStore) *UserFirendShipHandler {
	return &UserFirendShipHandler{
		userStore:   userStore,
		friendStore: friendStore,
		validate:    validator,
	}
}
