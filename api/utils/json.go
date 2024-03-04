package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/coderero/erochat-server/types"
	"github.com/labstack/echo/v4"
)

func ExtractInformation(err string) (string, string, string) {
	errMsg := fmt.Sprintf("%s ", err)
	castringError := regexp.MustCompile(`cannot unmarshal (.*?) into Go struct field (.*?) of type (.*?) `)

	if castringError.MatchString(errMsg) {
		field := strings.Split(castringError.FindStringSubmatch(errMsg)[2], ".")[1]
		givenType := castringError.FindStringSubmatch(errMsg)[1]
		expectedType := castringError.FindStringSubmatch(errMsg)[3]
		return field, givenType, expectedType
	}
	return "", "", ""
}

// Json Decode Function for JSON Decoding
func JSONDecode(c echo.Context, v interface{}) error {
	errRes := &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: "expected body type application/json",
	}
	if c.Request().Header.Get("Content-Type") != "application/json" {
		return errRes
	}

	if err := json.NewDecoder(c.Request().Body).Decode(v); err != nil {
		if errors.Is(err, io.EOF) {
			return errRes
		}
		errRes.Message = err.Error()
		return errRes
	}
	return nil
}

func JsonBindingErrorBuilder(err error) types.ApiResponse {
	field, givenType, expectedType := ExtractInformation(err.Error())
	errM := types.ApiResponse{
		Status:  types.Failure.String(),
		Code:    http.StatusBadRequest,
		Type:    types.ErrorTypeValidation.String(),
		Message: "invalid json",
		Errors: []types.Error{
			{
				Field:  field,
				Reason: fmt.Sprintf("expected %s, but got %s", expectedType, givenType),
			},
		},
	}
	return errM
}
