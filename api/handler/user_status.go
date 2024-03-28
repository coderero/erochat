package handler

import (
	"net/http"

	"github.com/coderero/erochat-server/api/utils"
	"github.com/coderero/erochat-server/interfaces"
	"github.com/coderero/erochat-server/types"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserStatusHandler struct {
	// userStore is a data store for user.
	userStore interfaces.UserStore

	// statusStore is a data store for status.
	statusStore interfaces.StatusStore

	// validate is a validator that validates the request.
	validate *validator.Validate
}

type CreateStatus struct {
	// Resource URL.
	URL string `json:"url" validate:"required"`

	// Resource Thumbnail.
	Thumbnail string `json:"thumbnail" validate:"required"`

	// Resource Title.
	Title string `json:"title" validate:"required"`
}

// NewUserStatusHandler returns a new user status handler.
func NewUserStatusHandler(validator *validator.Validate, userStore interfaces.UserStore, statusStore interfaces.StatusStore) *UserStatusHandler {
	return &UserStatusHandler{
		userStore:   userStore,
		statusStore: statusStore,
		validate:    validator,
	}
}

func (u *UserStatusHandler) GetStatus(c echo.Context) error {
	email, ok := c.Get("user").(string)
	if !ok {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to get user email",
		}
	}

	user, err := u.userStore.GetByEmail(email)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to get user",
		}
	}

	status, err := u.statusStore.GetStatus(user.ID)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to get status",
		}
	}

	return c.JSON(http.StatusOK, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "status retrieved successfully",
		Data:    status,
	})
}

func (u *UserStatusHandler) CreateStatus(c echo.Context) error {
	email, ok := c.Get("user").(string)
	if !ok {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to get user email",
		}
	}

	var status CreateStatus
	if err := utils.JSONDecode(c, &status); err != nil {
		return err
	}

	if err := u.validate.Struct(status); err != nil {
		return c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status:  types.Failure.String(),
			Code:    http.StatusBadRequest,
			Type:    types.ErrorTypeValidation.String(),
			Message: "validation error",
			Errors:  utils.ConvertValidationErrors(err),
		})
	}

	user, err := u.userStore.GetByEmail(email)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to get user",
		}
	}

	newStatus := &types.UserStatus{
		UserID:            user.ID,
		ResourceURI:       status.URL,
		ResourceThumbnail: status.Thumbnail,
		Title:             status.Title,
	}

	s, err := u.statusStore.CreateStatus(newStatus)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to create status",
		}
	}

	return c.JSON(http.StatusOK, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "status created successfully",
		Data:    s,
	})
}

func (u *UserStatusHandler) DeleteStatus(c echo.Context) error {
	statusID := c.Param("uid")
	if statusID == "" {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "status id is required",
		}
	}

	uid, err := uuid.Parse(statusID)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "invalid status id",
		}
	}

	email, ok := c.Get("user").(string)
	if !ok {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to get user email",
		}
	}

	user, err := u.userStore.GetByEmail(email)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to get user",
		}
	}

	err = u.statusStore.DeleteStatus(user.ID, uid)
	if err != nil {
		if err.Error() == "status not found" {
			return &echo.HTTPError{
				Code:    echo.ErrNotFound.Code,
				Message: "status not found",
			}
		}
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to delete status",
		}
	}

	return c.JSON(http.StatusOK, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "status deleted successfully",
	})

}
