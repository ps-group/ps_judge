package restapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Response - represents API response
type Response interface {
	// write - writes HTTP response and returns error that should be logged
	write(w http.ResponseWriter) error
}

// Ok - represets HttpStatusOk response
type Ok struct {
	Data interface{}
}

// InternalError - represents HttpInternalError response
type InternalError struct {
	Error error
}

// NewInternalError - creates internal error with reason string
func NewInternalError(reason string) *InternalError {
	return &InternalError{errors.New(reason)}
}

func (res *Ok) write(w http.ResponseWriter) error {
	bytes, err := json.Marshal(res.Data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := io.WriteString(w, err.Error())
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = io.WriteString(w, string(bytes))
	return err
}

func (res *InternalError) write(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, res.Error.Error())
	return res.Error
}
