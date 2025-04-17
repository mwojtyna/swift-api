package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

func ValidateStruct[T any](t T, validate *validator.Validate) error {
	err := validate.Struct(t)
	var ve validator.ValidationErrors

	if err != nil && errors.As(err, &ve) {
		msg := "Format checks failed for fields:\n"
		for _, fe := range ve {
			msg += fmt.Sprintf("'%s': %s\n", fe.Field(), fe.Tag())
		}

		return errors.New(msg)
	}

	return nil
}

func ReadJson[T any](w http.ResponseWriter, r *http.Request, t *T) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf(`"Content-Type" not "application/json"`)
	}

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		return err
	}

	return nil
}

func WriteJson[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteHttpError(w http.ResponseWriter, status int) {
	http.Error(w, "", status)
}

func NewApiServer(address string, db *sqlx.DB, logger *log.Logger) *ApiServer {
	validate := validator.New(validator.WithRequiredStructEnabled())
	// Return json name instead of struct name
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &ApiServer{
		address:  address,
		db:       db,
		logger:   logger,
		validate: validate,
	}
}

func (s *ApiServer) Run() {
	routerV1 := http.NewServeMux()
	routerV1.HandleFunc("GET /swift-codes/{swiftCode}", s.handleError(s.handleGetSwiftCodeV1))
	routerV1.HandleFunc("GET /swift-codes/country/{countryISO2code}", s.handleError(s.handleGetSwiftCodesForCountryV1))
	routerV1.HandleFunc("POST /swift-codes", s.handleError(s.handleAddSwiftCodeV1))
	routerV1.HandleFunc("DELETE /swift-codes/{swiftCode}", s.handleError(s.handleDeleteSwiftCodeV1))

	rootRouter := http.NewServeMux()
	rootRouter.Handle("/v1/", http.StripPrefix("/v1", routerV1))

	http.ListenAndServe(s.address, LoggingMiddleware(rootRouter, s.logger))
}

func (s *ApiServer) handleError(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteHttpError(w, http.StatusInternalServerError)
			s.logger.Printf(`ERROR on %s %s: "%s"`, r.Method, r.URL.Path, err)
		}
	})
}
