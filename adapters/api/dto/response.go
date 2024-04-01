package dto

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"golang-api-hexagonal/adapters/custom_error"
	"net/http"
)

// DefaultResponse create a default response object
func DefaultResponse(codeDescription, message string) map[string]interface{} {
	return map[string]interface{}{
		"code":    codeDescription,
		"message": message,
	}
}

// RenderResponse render http json response
func RenderResponse(ctx context.Context, writer http.ResponseWriter, httpStatusCode int, payload interface{}) {
	writer.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(ctx))
	writer.Header().Set("Content-Type", "application/json")

	marshal, err := json.Marshal(payload)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(httpStatusCode)
	_, _ = writer.Write(marshal)
}

// RenderErrorResponse render http error response
func RenderErrorResponse(ctx context.Context, writer http.ResponseWriter, httpStatusCode int, err error) {
	var response map[string]interface{}

	var customError *custom_error.StatusError
	if errors.As(err, &customError) {
		response = DefaultResponse(http.StatusText(customError.ErrorCode()), err.Error())
		RenderResponse(ctx, writer, customError.ErrorCode(), response)
		return
	}

	response = DefaultResponse(http.StatusText(httpStatusCode), err.Error())

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var details []map[string]interface{}
		for _, validationErr := range validationErrors {
			detail := map[string]interface{}{
				"field":       validationErr.Field(),
				"value":       validationErr.Value(),
				"location":    validationErr.Namespace(),
				"issue":       validationErr.Tag(),
				"description": validationErr.Error(),
			}
			details = append(details, detail)
		}
		response["message"] = "Validation errors"
		response["details"] = details
	}

	RenderResponse(ctx, writer, httpStatusCode, response)
}
