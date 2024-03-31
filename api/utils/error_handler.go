package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/coderero/erochat-server/types"
	"github.com/labstack/echo/v4"
)

// Custom HTTPErrorHandler is a custom error handler for HTTP errors.
func CustomHTTPErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var (
			errRes = types.ApiResponse{
				Status: types.Failure.String(),
			}
		)

		if c.Response().Committed {
			return
		}

		he, ok := err.(*echo.HTTPError)
		if ok {
			if he.Internal != nil {
				if herr, ok := he.Internal.(*echo.HTTPError); ok {
					he = herr
				}
			}
		} else {
			he = &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}

		}

		code := he.Code

		switch m := he.Message.(type) {
		case string:
			errRes.Code = code
			errRes.Type = extractErrorType(err).String()
			if code == echo.ErrInternalServerError.Code {
				errRes.Message = "internal server error"
			} else {
				errRes.Message = strings.ToLower(m)
			}
		case json.Marshaler:
		case error:
			errRes.Code = code
			errRes.Message = "internal server error"
			errRes.Type = types.ErrorTypeInternal.String()
		}

		// Send response
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(he.Code)
		} else {
			err = c.JSON(he.Code, errRes)
		}
		if err != nil {
			e.Logger.Error(err)
		}
	}
}

// Extract error type from error Message.
func extractErrorType(err error) types.ErrorType {
	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		return types.ErrorTypeUnknown
	}

	switch httpErr.Code {
	case http.StatusBadRequest:
		return types.ErrorBadRequest
	case http.StatusNotFound:
		return types.ErrorTypeNotFound
	case http.StatusConflict:
		return types.ErrorTypeConflict
	case http.StatusInternalServerError:
		return types.ErrorTypeInternal
	case http.StatusServiceUnavailable:
		return types.ErrorTypeServiceUnavailable
	case http.StatusUnauthorized:
		return types.ErrorTypeUnauthorized
	default:
		return types.ErrorTypeUnknown
	}
}
