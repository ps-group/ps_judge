package restapi

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func newRouter(config RouterConfig, context interface{}) *mux.Router {
	router := mux.NewRouter()
	subrouter := router.PathPrefix(config.APIPrefix).Subrouter()
	start := time.Now()

	decorateMethodHandler := func(handler MethodHandler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			request := requestImpl{
				request: r,
			}
			err := callMethodHandler(handler, context, &request, w)
			if err != nil {
				fields := logrus.Fields{
					"err":    err,
					"time":   time.Since(start),
					"method": r.Method,
					"url":    r.RequestURI,
				}
				logrus.WithFields(fields).Error("request failure")
			} else {
				fields := logrus.Fields{
					"time":   time.Since(start),
					"method": r.Method,
					"url":    r.RequestURI,
				}
				logrus.WithFields(fields).Info("done")
			}
		})
	}

	for _, route := range config.Routes {
		handler := decorateMethodHandler(route.Handler)
		subrouter.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(handler)
	}

	return router
}

// callMethodHandler - invoke handler, writes response, and stops panic if any happens.
func callMethodHandler(handler MethodHandler, context interface{}, r Request, w http.ResponseWriter) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("panic occured")
		}
	}()
	err = callMethodHandlerUnsafe(handler, context, r, w)
	return err
}

// callMethodHandlerUnsafe - invoke handler and writes response
func callMethodHandlerUnsafe(handler MethodHandler, context interface{}, r Request, w http.ResponseWriter) error {
	response := handler(context, r)
	return response.write(w)
}
