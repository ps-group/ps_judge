package restapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// Response - represents API response
// You must use predefined response implementation: Ok, BadRequest, etc
// If you initialize response with `error` interface, this error will be converted to JSON
//  using ErrorResponse structure serialization.
// Otherwise, any passed data will be converted to JSON directly.
type Response interface {
	// write - writes HTTP response and returns error that should be logged
	write(w http.ResponseWriter) error
}

// Ok - represets HttpStatusOk response
type Ok struct {
	Data interface{}
}

// BadRequest - represents HttpBadRequest response
type BadRequest struct {
	Data interface{}
}

// Unauthorized - represents HttpUnauthorized response
type Unauthorized struct {
	Data interface{}
}

// InternalError - represents HttpInternalError response
type InternalError struct {
	Data interface{}
}

// ErrorDetails - used internally to convert error into JSON
type ErrorDetails struct {
	Text string `json:"text"`
}

// ErrorResponse - used internally to convert error into JSON
type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

// NewInternalError - creates internal error with reason string
func NewInternalError(reason string) *InternalError {
	return &InternalError{errors.New(reason)}
}

func writeJSONResponse(data interface{}, status int, w http.ResponseWriter) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := io.WriteString(w, err.Error())
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_, err = io.WriteString(w, string(bytes))
	return err
}

func writeErrorResponse(err error, status int, w http.ResponseWriter) error {
	res := ErrorResponse{
		Error: ErrorDetails{
			Text: err.Error(),
		},
	}
	writeJSONResponse(res, status, w)
	return err
}

func writeResponse(data interface{}, status int, w http.ResponseWriter) error {
	if data == nil {
		w.WriteHeader(status)
		return nil
	}
	err, ok := data.(error)
	if ok {
		return writeErrorResponse(err, status, w)
	}
	return writeJSONResponse(data, status, w)
}

func (res *Ok) write(w http.ResponseWriter) error {
	return writeResponse(res.Data, http.StatusOK, w)
}

func (res *BadRequest) write(w http.ResponseWriter) error {
	return writeResponse(res.Data, http.StatusBadRequest, w)
}

func (res *Unauthorized) write(w http.ResponseWriter) error {
	return writeResponse(res.Data, http.StatusUnauthorized, w)
}

func (res *InternalError) write(w http.ResponseWriter) error {
	return writeResponse(res.Data, http.StatusInternalServerError, w)
}
