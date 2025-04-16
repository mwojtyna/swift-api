package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func NewAPIServer(address string, db *sqlx.DB, logger *log.Logger) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
		logger:  logger,
	}
}

func (s *APIServer) Run() {
	routerV1 := http.NewServeMux()
	routerV1.HandleFunc("GET /swift-codes/{swiftCode}", s.getSwiftCodeDetailsV1)
	routerV1.HandleFunc("GET /swift-codes/country/{countryISO2code}", s.getSwiftCodesForCountryV1)

	rootRouter := http.NewServeMux()
	rootRouter.Handle("/v1/", http.StripPrefix("/v1", routerV1))

	http.ListenAndServe(s.address, LoggingMiddleware(rootRouter, s.logger))
}
