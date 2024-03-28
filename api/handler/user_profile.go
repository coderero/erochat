package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/coderero/erochat-server/api/utils"
	"github.com/coderero/erochat-server/interfaces"
	"github.com/coderero/erochat-server/types"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	// validate is the validator.
	validate *validator.Validate

	// profileStore is a data store for profile.
	profileStore interfaces.ProfileStore

	// userStore is a data store for user.
	userStore interfaces.UserStore
}

// UserProfile is a user profile.
type UserProfile struct {
	// UID is the uuid of the profile.
	UID string `json:"uid"`

	// FirstName is the first name of the profile.
	FirstName string `json:"first_name"`

	// LastName is the last name of the profile.
	LastName string `json:"last_name"`

	// Bio is the bio of the profile.
	Bio string `json:"bio"`

	// Avatar is the avatar of the profile.
	Avatar string `json:"avatar"`

	// Username is the username of the profile.
	Username string `json:"username"`

	// Email is the email of the profile.
	Email string `json:"email"`

	// CreatedAt is the time the profile was created.
	CreatedAt time.Time `json:"created_at"`
}

type CreateProfile struct {
	// FirstName is the first name of the profile.
	FirstName string `json:"first_name" validate:"required"`

	// LastName is the last name of the profile.
	LastName string `json:"last_name" validate:"required"`

	// Bio is the bio of the profile.
	Bio string `json:"bio"`

	// Avatar is the avatar of the profile.
	Avatar string `json:"avatar" validate:"required"`
}

type UpdateProfile struct {
	// FirstName is the first name of the profile.
	FirstName string `json:"first_name"`

	// LastName is the last name of the profile.
	LastName string `json:"last_name"`

	// Bio is the bio of the profile.
	Bio string `json:"bio"`

	// Avatar is the avatar of the profile.
	Avatar string `json:"avatar"`
}

func NewProfileHandler(validator *validator.Validate, profileStore interfaces.ProfileStore, userStore interfaces.UserStore) *ProfileHandler {
	return &ProfileHandler{
		validate:     validator,
		profileStore: profileStore,
		userStore:    userStore,
	}
}

func (h *ProfileHandler) AddFriend(c echo.Context) error {
	uid := c.Param("uid")
	// Parse the uuid.
	id, err := uuid.Parse(uid)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "invalid profile id",
		}
	}

	email, ok := c.Get("user").(string)
	if !ok {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to get user email",
		}
	}

	addFriend := make(chan echo.Map, 2)
	go func(email string) {
		user, err := h.userStore.GetByEmail(email)
		if err != nil {
			addFriend <- echo.Map{
				"error": err,
			}
			return
		}

		addFriend <- echo.Map{
			"user": user,
		}
	}(email)

	go func(id uuid.UUID) {
		friend, err := h.profileStore.GetByUID(id)
		if err != nil {
			addFriend <- echo.Map{
				"error": err,
			}
			return
		}

		addFriend <- echo.Map{
			"friend": friend,
		}
	}(id)

	var (
		user      *types.User
		friend    *types.Profile
		userErr   error
		friendErr error
	)

	for i := 0; i < 2; i++ {
		res := <-addFriend
		if res["user"] != nil {
			user = res["user"].(*types.User)
		} else if res["friend"] != nil {
			friend = res["friend"].(*types.Profile)
		} else if res["error"] != nil {
			if userErr == nil {
				userErr = res["error"].(error)
			} else {
				friendErr = res["error"].(error)
			}
		}
	}

	if userErr != nil || friendErr != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "failed to get user or friend",
		}
	}

	if user.ID == friend.UserID {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "cannot add yourself as a friend",
		}
	}

	err = h.profileStore.CreateFriendship(user.UID.String(), friend.UID.String())
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "failed to create friendship",
		}
	}

	return c.JSON(http.StatusOK, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "friend request sent successfully",
	})

}

// GetProfileByID returns a profile by its uuid.
func (h *ProfileHandler) GetProfileByID(c echo.Context) error {
	uid := c.Param("uid")

	// Parse the uuid.
	id, err := uuid.Parse(uid)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "invalid profile id",
		}
	}

	profile, err := h.profileStore.GetByUID(id)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: "profile not found",
		}
	}

	profileResponse := UserProfile{
		UID:       profile.UID.String(),
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Bio:       profile.Bio,
		Avatar:    profile.Avatar,
		Username:  profile.Username,
		Email:     profile.Email,
		CreatedAt: profile.CreatedAt,
	}

	return c.JSON(http.StatusOK, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "profile retrieved successfully",
		Data:    profileResponse,
	})
}

// GetProfileByID returns a profile by its uuid.
func (h *ProfileHandler) GetProfile(c echo.Context) error {
	email, ok := c.Get("user").(string)
	if !ok {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to get user email",
		}
	}

	profile, err := h.profileStore.GetByEmail(email)
	if err != nil {
		return err
	}

	profileResponse := UserProfile{
		UID:       profile.UID.String(),
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Bio:       profile.Bio,
		Avatar:    profile.Avatar,
		Username:  profile.Username,
		Email:     profile.Email,
		CreatedAt: profile.CreatedAt,
	}

	return c.JSON(http.StatusOK, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "profile retrieved successfully",
		Data:    profileResponse,
	})
}

// CreateProfile creates a new profile.
func (h *ProfileHandler) CreateProfile(c echo.Context) error {
	profile := new(CreateProfile)
	if err := utils.JSONDecode(c, profile); err != nil {
		if strings.Contains(err.Error(), "json:") {
			return c.JSON(http.StatusBadRequest, utils.JsonBindingErrorBuilder(err))
		}
		return err
	}

	// Validate the request body.
	if err := h.validate.Struct(profile); err != nil {
		return c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status:  types.Failure.String(),
			Code:    http.StatusBadRequest,
			Message: types.ErrorTypeValidation.String(),
			Errors:  utils.ConvertValidationErrors(err),
		})
	}

	email, ok := c.Get("user").(string)
	if !ok {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to get user email",
		}
	}

	// Get the user by its email.
	user, err := h.userStore.GetByEmail(email)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrUnauthorized.Code,
			Message: "unauthorized",
		}
	}

	profileData := &types.Profile{
		UID:       user.UID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Bio:       profile.Bio,
		Avatar:    profile.Avatar,
		UserID:    user.ID,
	}

	res, err := h.profileStore.Create(profileData)
	if err != nil {
		if err == interfaces.ErrProfileExists {
			return &echo.HTTPError{
				Code:    http.StatusConflict,
				Message: "profile already exists",
			}
		}
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	profileResponse := UserProfile{
		UID:       res.UID.String(),
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Avatar:    res.Avatar,
		Bio:       res.Bio,
		Username:  res.Username,
		Email:     res.Email,
		CreatedAt: res.CreatedAt,
	}

	return c.JSON(http.StatusCreated, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusCreated,
		Message: "profile created successfully",
		Data:    profileResponse,
	})
}

// UpdateProfile updates a profile.
func (h *ProfileHandler) UpdateProfile(c echo.Context) error {
	profile := new(UpdateProfile)
	if err := utils.JSONDecode(c, profile); err != nil {
		if strings.Contains(err.Error(), "json:") {
			return c.JSON(http.StatusBadRequest, utils.JsonBindingErrorBuilder(err))
		}
		return err
	}

	email, ok := c.Get("user").(string)
	if !ok {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to get user email",
		}
	}

	if profile.FirstName == "" && profile.LastName == "" && profile.Avatar == "" {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "at least one field is required",
		}
	}

	// Get the user by its email.
	user, err := h.userStore.GetByEmail(email)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrUnauthorized.Code,
			Message: "unauthorized",
		}
	}

	profileData := &types.Profile{
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Bio:       profile.Bio,
		Avatar:    profile.Avatar,
		UserID:    user.ID,
	}

	res, err := h.profileStore.Update(profileData)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	profileResponse := UserProfile{
		UID:       res.UID.String(),
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Avatar:    res.Avatar,
		Username:  res.Username,
		Email:     res.Email,
		CreatedAt: res.CreatedAt,
	}

	return c.JSON(http.StatusOK, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "profile updated successfully",
		Data:    profileResponse,
	})
}

// DeleteProfile deletes a profile by its uuid.
func (h *ProfileHandler) DeleteProfile(c echo.Context) error {
	email, ok := c.Get("user").(string)
	if !ok {
		return &echo.HTTPError{
			Code:    echo.ErrBadRequest.Code,
			Message: "failed to get user email",
		}
	}

	// Get the user by its email.
	user, err := h.userStore.GetByEmail(email)
	if err != nil {
		return &echo.HTTPError{
			Code:    echo.ErrUnauthorized.Code,
			Message: "unauthorized",
		}
	}

	_, err = h.profileStore.Delete(user.ID)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(http.StatusOK, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "profile deleted successfully",
	})
}
