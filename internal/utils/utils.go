package utils

import (
	"encoding/json"
	"errors"
	"github.com/dominik24c/quizzes_api/internal"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
)

func ErrorHandler(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	resp := make(map[string]string)
	resp["message"] = message
	json.NewEncoder(w).Encode(resp)
}

func JsonResponseHandler(w http.ResponseWriter, status int, message interface{}) {
	if status != 0 {
		w.WriteHeader(status)
	}
	json.NewEncoder(w).Encode(message)
}

func IsAuthorized(authHeader string) (internal.AuthResponse, error) {
	var authData internal.AuthResponse

	c := &http.Client{}
	req, _ := http.NewRequest("POST", "http://auth_service:9998/auth", nil)
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/json")
	response, err := c.Do(req)

	if err != nil {
		return authData, err
	}

	err = json.NewDecoder(response.Body).Decode(&authData)
	if err != nil {
		return authData, err
	}

	return authData, nil
}

func Pagination(r *http.Request, limit int64) *options.FindOptions {
	query := r.URL.Query()
	page := int64(1)
	if p := query.Get("page"); p != "" {
		pageTmp, err := strconv.ParseInt(p, 10, 64)
		if err == nil && pageTmp > 0 {
			page = pageTmp
		}
	}

	opts := options.FindOptions{}
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	return &opts
}

func ValidateRequestBody(s interface{}) (*internal.ErrorMessage, error) {
	var errorMsg *internal.ErrorMessage
	v := validator.New()
	if err := v.Struct(s); err != nil {

		validationErrors := err.(validator.ValidationErrors)
		var errorsData []internal.ErrorField
		for _, validationError := range validationErrors {
			errorsData = append(errorsData, internal.ErrorField{
				Field: validationError.Field(),
				Value: validationError.Tag(),
			})
		}

		errorMsg = &internal.ErrorMessage{Errors: errorsData}
		return errorMsg, errors.New("invalid input")
	}

	return errorMsg, nil
}
