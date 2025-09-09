package dto

import (
	"encoding/json"
	"net/http"

	"github.com/sagungw/gotrunks/errors"
)

var (
	errorHttpStatusMapping = map[errors.Code]int{
		errors.CodeInternalError:       http.StatusInternalServerError,
		errors.CodeAuthenticationError: http.StatusUnauthorized,
		errors.CodeParameterError:      http.StatusBadRequest,
	}

	statusFailure = "FAILURE"
	statusSuccess = "SUCCESS"
)

type Response struct {
	Status string         `json:"status"`
	Error  *ErrorResponse `json:"error,omitempty"`
	Data   any            `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	SubCode string `json:"subcode"`
	Message string `json:"message"`
}

func WriteJson(w http.ResponseWriter, httpStatus int, data any) {
	response := &Response{
		Status: statusSuccess,
		Data:   data,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func WriteJsonError(w http.ResponseWriter, err error) {
	WriteJsonErrorData(w, err, nil)
}

func WriteJsonErrorData(w http.ResponseWriter, err error, data any) {
	response := &Response{
		Status: statusFailure,
		Error:  NewErrorResponse(err),
		Data:   data,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(mapErrorToHttpStatus(err))
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NewErrorResponse(err error) *ErrorResponse {
	e, ok := err.(*errors.CodedError)
	if !ok {
		e = errors.NewCodedError(errors.CodeInternalError, errors.SubCodeNone, err)
	}

	return &ErrorResponse{
		Code:    string(e.Code),
		SubCode: string(e.SubCode),
		Message: e.Unwrap().Error(),
	}
}

func mapErrorToHttpStatus(err error) int {
	e, ok := err.(*errors.CodedError)
	if !ok {
		return http.StatusInternalServerError
	}

	return errorHttpStatusMapping[e.Code]
}
