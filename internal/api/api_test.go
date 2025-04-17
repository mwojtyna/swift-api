package api

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateStruct(t *testing.T) {
	type ts struct {
		Field1 string `validate:"required"`
	}

	tests := []struct {
		name    string
		given   ts
		wantErr bool
	}{
		{"validated", ts{Field1: "abc"}, false},
		{"validation error", ts{}, true},
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.given, validate)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReadJson(t *testing.T) {
	type ts struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name           string
		contentType    string
		requestBody    string
		target         *ts
		expectedResult ts
		wantErr        bool
	}{
		{
			name:           "parsed",
			contentType:    "application/json",
			requestBody:    `{"name":"John","age":30}`,
			target:         &ts{},
			expectedResult: ts{Name: "John", Age: 30},
		},
		{
			name:        "wrong content type",
			contentType: "text/plain",
			requestBody: `{"name":"John","age":30}`,
			target:      &ts{},
			wantErr:     true,
		},
		{
			name:        "malformed JSON",
			contentType: "application/json",
			requestBody: `{"name":"John","age":30`,
			target:      &ts{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			err := ReadJson(w, req, tt.target)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, *tt.target)
			}
		})
	}
}

func TestWriteJson(t *testing.T) {
	type ts struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name       string
		statusCode int
		data       ts
		wantErr    bool
	}{
		{
			name:       "parsed",
			statusCode: 418,
			data:       ts{Name: "John", Age: 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := WriteJson(w, tt.statusCode, tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
				assert.Equal(t, tt.statusCode, w.Code)
			}
		})
	}
}

func TestWriteHttpError(t *testing.T) {
	tests := []struct {
		statusCode int
	}{
		{200},
		{418},
	}

	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.statusCode), func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteHttpError(w, tt.statusCode)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestHandleError(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name           string
		handlerFunc    func(http.ResponseWriter, *http.Request) error
		expectedStatus int
		logContains    string
	}{
		{
			name: "handler returns no error",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) error {
				w.WriteHeader(http.StatusOK)
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "handler returns error",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("test error")
			},
			expectedStatus: http.StatusInternalServerError,
			logContains:    "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup logger with buffer to capture output
			var logBuf bytes.Buffer
			logger := log.New(&logBuf, "", 0)

			server := &ApiServer{
				logger:   logger,
				validate: validator.New(),
				// db and address not needed
			}

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			wrappedHandler := server.handleError(tt.handlerFunc)
			wrappedHandler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Result().StatusCode)

			logOutput := logBuf.String()
			if tt.logContains != "" {
				assert.Contains(t, logOutput, tt.logContains)
			} else {
				assert.Empty(t, logOutput)
			}
		})
	}
}
