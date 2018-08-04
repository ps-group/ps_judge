package restapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// Request - provides access to request parameters and other information for the API handler
type Request interface {
	Var(name string) string
	ReadJSON(result interface{}) error
}

// requestImpl - wrapper for http.Request which implements Request interface
type requestImpl struct {
	request *http.Request
	vars    map[string]string
}

// Var - parses argument from URL (for GET)
func (req *requestImpl) Var(name string) string {
	if req.vars == nil {
		req.vars = mux.Vars(req.request)
	}
	return req.vars[name]
}

// ReadJSON - read JSON from request body or returns error
func (req *requestImpl) ReadJSON(result interface{}) error {
	bytes, err := ioutil.ReadAll(req.request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, result)
	if err != nil {
		return err
	}
	return nil
}
