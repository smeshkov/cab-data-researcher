package app

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/smeshkov/cab-data-researcher/app/handlers"
	"github.com/smeshkov/cab-data-researcher/cfg"
	"github.com/smeshkov/cab-data-researcher/cache"
)

// CreateHandler creates handler for application.
func CreateHandler(env string, config *cfg.Config, cdb cache.CabDBCache) http.Handler {
	// Use gorilla/mux for rich routing.
	// See http://www.gorillatoolkit.org/pkg/mux
	r := mux.NewRouter()

	// Indicate the server is healthy.
	r.Methods("GET").Path("/health").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

	// ------------------------ API ------------------------
	api := r.PathPrefix("/api/v1").Subrouter()

	api.Methods(http.MethodPost).
		Path("/trip/count").
		Handler(handlers.AppHandler(getTripCount(cdb)))

	api.Methods(http.MethodPost).
		Path("/cache/clear").
		Handler(handlers.AppHandler(clearCache(cdb)))

	return r
}
