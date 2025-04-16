package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func ReadJSON[T any](w http.ResponseWriter, r *http.Request, t *T) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf(`"Content-Type" not "application/json"`)
	}

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func WriteHTTPError(w http.ResponseWriter, status int) {
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
	routerV1.HandleFunc("GET /swift-codes/{swiftCode}", s.handleError(s.handleGetSwiftCodeV1))
	routerV1.HandleFunc("GET /swift-codes/country/{countryISO2code}", s.handleError(s.handleGetSwiftCodesForCountryV1))
	routerV1.HandleFunc("POST /swift-codes", s.handleError(s.handleAddSwiftCodeV1))

	rootRouter := http.NewServeMux()
	rootRouter.Handle("/v1/", http.StripPrefix("/v1", routerV1))

	http.ListenAndServe(s.address, LoggingMiddleware(rootRouter, s.logger))
}

func (s *APIServer) handleError(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteHTTPError(w, http.StatusInternalServerError)
			s.logger.Printf(`ERROR on %s %s: "%s"`, r.Method, r.URL.Path, err)
		}
	})
}
