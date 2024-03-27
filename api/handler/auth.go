package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/coderero/erochat-server/api/utils"
	"github.com/coderero/erochat-server/interfaces"
	"github.com/coderero/erochat-server/types"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// AuthHandler represents an HTTP handler for authentication.
type AuthHandler struct {
	// Validator for request validation.
	validator *validator.Validate

	// UserStore represents a user store.
	userStore interfaces.UserStore

	// PasswordHasher represents a password hasher.
	passwordHasher interfaces.PassService

	// TokenService represents a token service.
	tokenService interfaces.TokenService
}

// AuthCreate represents a request to create a new user.
type AuthCreate struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Auth represents a request to login a user.
type Auth struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshToken represents a request to refresh a token.
type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(userStore interfaces.UserStore, passwordHasher interfaces.PassService, tokenService interfaces.TokenService) *AuthHandler {
	return &AuthHandler{
		validator:      validator.New(),
		userStore:      userStore,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
	}
}

var (
	invalidCred = types.ApiResponse{
		Status:  types.Failure.String(),
		Code:    http.StatusBadRequest,
		Type:    types.ErrorTypeInvalidCredentials.String(),
		Message: "invalid credentials provided",
	}
	validationErr = types.ApiResponse{
		Status:  types.Failure.String(),
		Code:    http.StatusBadRequest,
		Type:    types.ErrorTypeValidation.String(),
		Message: "validation error",
	}
)

// Login logs in a user.
func (h *AuthHandler) Login(c echo.Context) error {
	var (
		params Auth
		user   *types.User
		err    error
	)
	if err := utils.JSONDecode(c, &params); err != nil {
		if strings.Contains(err.Error(), "json:") {
			return c.JSON(http.StatusBadRequest, utils.JsonBindingErrorBuilder(err))
		}
	}

	if err := checkForLoginParams(params); len(err) > 0 {
		validationErr.Errors = err
		return c.JSON(http.StatusBadRequest, validationErr)
	}

	// Check if the user exists.
	if params.Username != "" {
		user, err = h.userStore.GetByUsername(params.Username)
	} else {
		user, err = h.userStore.GetByEmail(params.Email)
	}

	// If the user is not found, return an error.
	if err != nil {
		if errors.Is(err, interfaces.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, types.ApiResponse{
				Status:  types.Failure.String(),
				Code:    http.StatusNotFound,
				Type:    types.ErrorTypeNotFound.String(),
				Message: "user not found",
			})
		}
		return c.JSON(http.StatusBadRequest, userStoreErrResBuilder(err))
	}

	if user.DeletedAt.Valid {
		// Send a response that your account has been deleted and you can't login
		// You can also send a link to recover the account
		return c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status:  types.Failure.String(),
			Code:    http.StatusBadRequest,
			Type:    types.ErrorTypeAccountDeleted.String(),
			Message: "your account has been deleted, although you can recover it with the link sent to your email",
		})
	}

	// Check if the password is valid.
	if !h.passwordHasher.Compare(params.Password, user.Password) {
		return c.JSON(http.StatusBadRequest, invalidCred)
	}

	// Generate a token and a refresh token.
	token, refreshToken, err := h.tokenService.GenerateTokens(user.Email, user.UID)

	// If an error occurred, return it.
	if err != nil {
		return c.JSON(http.StatusBadRequest, userStoreErrResBuilder(err))
	}

	// Save the tokens in the cookies.
	utils.SaveCookie(c, "__a", token)
	utils.SaveCookie(c, "__r", refreshToken)

	return c.JSON(http.StatusOK, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "user created successfully",
		Data: echo.Map{
			"access_token":  token,
			"refresh_token": refreshToken,
		},
	})
}

// Register registers a new user.
func (h *AuthHandler) Register(c echo.Context) error {
	var (
		params AuthCreate
		user   *types.User
		err    error
	)
	if err := utils.JSONDecode(c, &params); err != nil {
		if strings.Contains(err.Error(), "json:") {
			return c.JSON(http.StatusBadRequest, utils.JsonBindingErrorBuilder(err))
		}
		return err
	}

	// Validate the request.
	if err := h.validator.Struct(params); err != nil {
		validationErr.Errors = utils.ConvertValidationErrors(err)
		return c.JSON(http.StatusBadRequest, validationErr)
	}

	// create a channel to check if the username or email exists both concurrently
	exits := make(chan map[string]interface{}, 2)

	// Check if the username exists.
	go func() {
		user, _ = h.userStore.GetByUsername(params.Username)
		if user != nil {
			if user.DeletedAt.Valid {
				exits <- map[string]interface{}{"username": false, "deleted": true}
			}
			exits <- map[string]interface{}{"username": true}
		} else {
			exits <- map[string]interface{}{"username": false}
		}
	}()

	// Check if the email exists.
	go func() {
		user, _ = h.userStore.GetByEmail(params.Email)
		if user != nil {
			if user.DeletedAt.Valid {
				exits <- map[string]interface{}{"email": false, "deleted": true}
			}
			exits <- map[string]interface{}{"email": true}
		} else {
			exits <- map[string]interface{}{"email": false}
		}
	}()

	// Create a list of the errors that occurred.
	var errors []types.Error

	// Check if the username or email exists.
	for i := 0; i < 2; i++ {
		exist := <-exits
		for k, v := range exist {
			if v.(bool) {
				errors = append(errors, types.Error{
					Field:  k,
					Reason: k + " already exists",
				})
			}
		}
	}

	// If there are errors, return them.
	if len(errors) > 0 {
		validationErr.Errors = errors
		return c.JSON(http.StatusBadRequest, validationErr)
	}

	// Hash the password.
	hashedPass, err := h.passwordHasher.Hash(params.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, userStoreErrResBuilder(err))
	}

	// Create the user.
	user = &types.User{
		Username: params.Username,
		Email:    params.Email,
		Password: hashedPass,
	}

	// Store the user.
	user, err = h.userStore.Create(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, userStoreErrResBuilder(err))
	}

	// Generate a token and a refresh token.
	token, refreshToken, err := h.tokenService.GenerateTokens(user.Email, user.UID)

	// If an error occurred, return it.
	if err != nil {
		return c.JSON(http.StatusBadRequest, userStoreErrResBuilder(err))
	}

	// Save the tokens in the cookies.
	utils.SaveCookie(c, "__a", token)
	utils.SaveCookie(c, "__r", refreshToken)

	return c.JSON(http.StatusOK, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "user created successfully",
		Data: echo.Map{
			"access_token":  token,
			"refresh_token": refreshToken,
		},
	})

}

// ! Non Web Handler
// RefreshToken refreshes a token.
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var (
		refreshToken RefreshToken
		err          error
	)
	// Get the access token from the request Header.
	bearerToken := c.Request().Header.Get("Authorization")
	if bearerToken == "" {
		return c.JSON(http.StatusUnauthorized, types.ApiResponse{
			Status:  types.Failure.String(),
			Code:    http.StatusUnauthorized,
			Type:    types.ErrorTypeUnauthorized.String(),
			Message: "unauthorized",
		})
	}

	_ = strings.Split(bearerToken, "Bearer ")[1]

	if err := utils.JSONDecode(c, &refreshToken); err != nil {
		if strings.Contains(err.Error(), "json:") {
			return c.JSON(http.StatusBadRequest, utils.JsonBindingErrorBuilder(err))
		}
	}

	// Validate the request.
	if err := h.validator.Struct(refreshToken); err != nil {
		validationErr.Errors = utils.ConvertValidationErrors(err)
		return c.JSON(http.StatusBadRequest, validationErr)
	}

	// TODO: Revoking the existing access token.

	// Check if the refresh token is valid.
	token, err := h.tokenService.RefreshToken(refreshToken.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, userStoreErrResBuilder(err))
	}

	// Return the new token.
	return c.JSON(http.StatusOK, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "token refreshed successfully",
		Data:    token,
	})
}

// Logout logs out a user.
func (h *AuthHandler) Logout(c echo.Context) error {
	// Get the access token and the refresh token from the cookies.
	accessToken := utils.GetCookie(c, "__a")
	refreshToken := utils.GetCookie(c, "__r")

	if accessToken == "" || refreshToken == "" {
		return c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status:  types.Failure.String(),
			Code:    http.StatusBadRequest,
			Type:    types.ErrorInvalidRequest.String(),
			Message: "you are not logged in",
		})
	}

	// Delete the cookies.
	utils.DeleteCookie(c, "__a")
	utils.DeleteCookie(c, "__r")
	return c.JSON(http.StatusOK, types.ApiResponse{
		Status:  types.Success.String(),
		Code:    http.StatusOK,
		Message: "user logged out successfully",
	})
}

// checkForLoginParams checks if the login parameters are valid.
func checkForLoginParams(a Auth) []types.Error {
	var (
		errors []types.Error
	)
	if a.Username == "" && a.Email == "" {
		errors = append(errors, types.Error{
			Field:  "username",
			Reason: "username or email is required",
		})
		errors = append(errors, types.Error{
			Field:  "email",
			Reason: "username or email is required",
		})
	}

	if a.Username != "" && a.Email != "" {
		errors = append(errors, types.Error{
			Field:  "username",
			Reason: "username and email cannot be used together",
		})
		errors = append(errors, types.Error{
			Field:  "email",
			Reason: "username and email cannot be used together",
		})
	}

	if a.Email != "" && !utils.IsValidEmail(a.Email) {
		errors = append(errors, types.Error{
			Field:  "email",
			Reason: "email is invalid",
		})
	}

	if a.Password == "" {
		errors = append(errors, types.Error{
			Field:  "password",
			Reason: "field is required",
		})
	}

	if a.Password != "" && len(a.Password) < 8 {
		errors = append(errors, types.Error{
			Field:  "password",
			Reason: "password must be at least 8 characters",
		})
	}
	return errors
}

func userStoreErrResBuilder(err error) types.ApiResponse {
	var (
		apiRes types.ApiResponse = types.ApiResponse{
			Status: types.Failure.String(),
			Code:   http.StatusBadRequest,
		}
		errors []types.Error
	)

	switch err {
	case interfaces.ErrUserNotFound:
		apiRes.Code = http.StatusNotFound
		apiRes.Type = types.ErrorTypeNotFound.String()
	case interfaces.ErrFailedToGetUser:
		apiRes.Code = http.StatusBadRequest
		apiRes.Type = types.ErrorTypeInternal.String()
		apiRes.Message = "failed to get user"
	case interfaces.ErrFailedToCreateUser:
		apiRes.Code = http.StatusBadRequest
		apiRes.Type = types.ErrorTypeInternal.String()
		apiRes.Message = "failed to create user"
	case interfaces.ErrFailedToUpdateUser:
		apiRes.Code = http.StatusBadRequest
		apiRes.Type = types.ErrorTypeInternal.String()
		apiRes.Message = "failed to update user"
	case interfaces.ErrFailedToDeleteUser:
		apiRes.Code = http.StatusBadRequest
		apiRes.Type = types.ErrorTypeInternal.String()
		apiRes.Message = "failed to delete user"
	case interfaces.ErrEmailExists:
		apiRes.Code = http.StatusConflict
		apiRes.Type = types.ErrorTypeConflict.String()
		apiRes.Message = "conflict occurred"
		errors = append(errors, types.Error{
			Field:  "email",
			Reason: "email already exists",
		})
		apiRes.Errors = append(apiRes.Errors, errors...)
	case interfaces.ErrUsernameExists:
		apiRes.Code = http.StatusConflict
		apiRes.Type = types.ErrorTypeConflict.String()
		apiRes.Message = "conflict occurred"
		errors = append(errors, types.Error{
			Field:  "username",
			Reason: "username already exists",
		})
		apiRes.Errors = append(apiRes.Errors, errors...)
	default:
		apiRes.Code = http.StatusBadRequest
		apiRes.Type = types.ErrorTypeUnknown.String()
		apiRes.Message = "unknown error occurred"
	}
	return apiRes
}
