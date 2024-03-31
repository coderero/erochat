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
	ResourceURL string `json:"resource_url" validate:"required"`

	// Resource Thumbnail.
	ResourceThumbnail string `json:"resource_thumbnail" validate:"required"`

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
	r_uid, ok := c.Get("uid").(string)
	if !ok {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "something went wrong",
		}
	}

	uid, err := uuid.Parse(r_uid)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "something went wrong",
		}
	}
	status, err := u.statusStore.GetStatus(uid)
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
	r_uid, ok := c.Get("uid").(string)
	if !ok {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "something went wrong",
		}
	}

	uid, err := uuid.Parse(r_uid)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "something went wrong",
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

	newStatus := &types.UserStatus{
		UserID:            uid,
		ResourceURI:       status.ResourceURL,
		ResourceThumbnail: status.ResourceThumbnail,
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

	r_uid, ok := c.Get("uid").(string)
	if !ok {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "something went wrong",
		}
	}

	user_uid, err := uuid.Parse(r_uid)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "something went wrong",
		}
	}

	err = u.statusStore.DeleteStatus(user_uid, uid)
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
