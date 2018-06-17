package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// APIContext - provides access to request parameters and other information for the API handler
type APIContext interface {
	Vars() map[string]string
	ReadJSON(result interface{}) error
	WriteJSON(result interface{}) error
	ConnectDB() (*sql.DB, error)
}

// APIHandler - handler which handles request and returns serializable interface or an error
type APIHandler func(ctx APIContext) error

// RawAPIHandler - handler which directly works with request and response
type RawAPIHandler func(w http.ResponseWriter, r *http.Request) error

// DefaultAPIContext - provides default JSON-based RESTful API context
type DefaultAPIContext struct {
	request   *http.Request
	writer    http.ResponseWriter
	connector DatabaseConnector
}

// ParseArgs - parses arguments from JSON (for POST/PUT/DELETE) or from URL (for GET)
func (c *DefaultAPIContext) Vars() map[string]string {
	return mux.Vars(c.request)
}

func (c *DefaultAPIContext) WriteJSON(result interface{}) error {
	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	c.writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = io.WriteString(c.writer, string(bytes))
	if err != nil {
		return err
	}
	return nil
}

func (c *DefaultAPIContext) ReadJSON(result interface{}) error {
	bytes, err := ioutil.ReadAll(c.request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, result)
	if err != nil {
		return err
	}
	return nil
}

func (c *DefaultAPIContext) ConnectDB() (*sql.DB, error) {
	return c.connector.Connect()
}

func openFileLogger(filename string) *os.File {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(file)
	return file
}

func newRouter(connector DatabaseConnector) *mux.Router {
	router := mux.NewRouter()
	subrouter := router.PathPrefix(BuilderAPIPrefix).Subrouter()
	start := time.Now()

	decorateWithLog := func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			inner.ServeHTTP(w, r)

			fields := logrus.Fields{
				"time":   time.Since(start),
				"method": r.Method,
				"url":    r.RequestURI,
			}
			logrus.WithFields(fields).Info("done")
		})
	}

	decorateWithJSON := func(inner APIHandler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapper := DefaultAPIContext{
				request:   r,
				writer:    w,
				connector: connector,
			}
			err := inner(&wrapper)
			if err != nil {
				fields := logrus.Fields{
					"err":    err,
					"time":   time.Since(start),
					"method": r.Method,
					"url":    r.RequestURI,
				}
				logrus.WithFields(fields).Error("request failure")
				io.WriteString(w, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		})
	}

	for _, route := range jsonRoutes {
		handler := decorateWithLog(decorateWithJSON(route.HandlerFunc))
		subrouter.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(handler)
	}

	return router
}
