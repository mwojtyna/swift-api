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
	http.HandleFunc("/v1/swift-codes/{swiftCode}", s.getSwiftCodeDetailsV1)

	http.ListenAndServe(s.address, nil)
}
