package restapi

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
			response := handler(context, &request)
			err := response.write(w)
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
