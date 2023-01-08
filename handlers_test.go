package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"promhsd/db"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	storage testStorage
)

type testStorage struct {
	returnError error
	//returnItem  *Target
}

func (s *testStorage) Create(*db.Target) error {

	return s.returnError
}

func (s *testStorage) Update(*db.Target) error {
	return s.returnError
}

func (s *testStorage) Delete(*db.Target) error {
	return s.returnError
}

func (s *testStorage) Get(*db.Target) error {
	return s.returnError
}

func (s *testStorage) GetAll(*[]db.Target) error {
	return s.returnError
}

func (s *testStorage) IsHealthy() bool {
	return true
}

type testStorageService struct{}

func (s *testStorageService) ServiceID() string {
	return "testdb"
}

func (s *testStorageService) New(opt string) (db.Storage, error) {
	return &storage, nil
}

func TestMain(m *testing.M) {
	db.RegisterStorage(&testStorageService{})

	os.Exit(m.Run())
}

func Test_getTargetsHandler(t *testing.T) {
	var (
		err error
	)
	dbService, err = db.New("testdb", "")
	assert.NoError(t, err)

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/targets/", nil)
	storage.returnError = &db.StorageError{}
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	storage.returnError = nil
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func Test_getTargetHandler(t *testing.T) {
	var (
		err error
	)
	dbService, err = db.New("testdb", "")
	assert.NoError(t, err)

	router := setupRouter()

	tests := []struct {
		name string
		err  error
		code int
	}{
		{
			name: "error500",
			err:  &db.StorageError{},
			code: http.StatusInternalServerError,
		},
		{
			name: "NoErrors",
			err:  nil,
			code: http.StatusOK,
		},
		{
			name: "validationError",
			err:  &db.ValidationError{},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "notFoundError",
			err:  &db.NotFoundError{},
			code: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/target/1", nil)
			storage.returnError = tt.err
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.code, w.Code)
		})
	}
}

func Test_createTargetHandler(t *testing.T) {
	var (
		err error
	)
	dbService, err = db.New("testdb", "")
	assert.NoError(t, err)

	router := setupRouter()

	tests := []struct {
		name    string
		err     error
		code    int
		payload string
	}{
		{
			name:    "error400",
			err:     &db.StorageError{},
			code:    http.StatusBadRequest,
			payload: "",
		},
		{
			name:    "error500",
			err:     &db.StorageError{},
			code:    http.StatusInternalServerError,
			payload: `{"name": "test", "entries": [{"targets": "127.0.0.1:5000", "labels": "key=val,k1=v1"}]}`,
		},
		{
			name:    "NoErrors",
			err:     nil,
			code:    http.StatusOK,
			payload: `{"name": "test", "entries": [{"targets": "127.0.0.1:5000", "labels": "key=val,k1=v1"}]}`,
		},
		{
			name:    "validationError",
			err:     &db.ValidationError{},
			code:    http.StatusUnprocessableEntity,
			payload: `{"name": "test", "entries": [{"targets": "127.0.0.1:5000", "labels": "key=val,k1=v1"}]}`,
		},
		{
			name:    "nameEmpty",
			err:     nil,
			code:    http.StatusBadRequest,
			payload: `{"name": "", "entries": [{"targets": "127.0.0.1:5000", "labels": "key=val,k1=v1"}]}`,
		},
		{
			name:    "validationError2",
			err:     nil,
			code:    http.StatusUnprocessableEntity,
			payload: `{"name": "test", "entries": [{"targets": "127.0.0.1:5000", "labels": "keys"}]}`,
		},
		{
			name:    "conflictError",
			err:     &db.ConflictError{},
			code:    http.StatusConflict,
			payload: `{"name": "test", "entries": [{"targets": "127.0.0.1:5000", "labels": "key=val,k1=v1"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/target/", strings.NewReader(tt.payload))
			storage.returnError = tt.err
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.code, w.Code)
		})
	}
}

func Test_updateTargetHandler(t *testing.T) {
	var (
		err error
	)
	dbService, err = db.New("testdb", "")
	assert.NoError(t, err)

	router := setupRouter()

	tests := []struct {
		name    string
		err     error
		code    int
		payload string
	}{
		{
			name:    "error400",
			err:     &db.StorageError{},
			code:    http.StatusBadRequest,
			payload: "",
		},
		{
			name:    "error500",
			err:     &db.StorageError{},
			code:    http.StatusInternalServerError,
			payload: `{"name": "test", "entries": [{"targets": "127.0.0.1:5000", "labels": "key=val,k1=v1"}]}`,
		},
		{
			name:    "NoErrors",
			err:     nil,
			code:    http.StatusOK,
			payload: `{"name": "test", "entries": [{"targets": "127.0.0.1:5000", "labels": "key=val,k1=v1"}]}`,
		},
		{
			name:    "validationError",
			err:     &db.ValidationError{},
			code:    http.StatusUnprocessableEntity,
			payload: `{"name": "test", "entries": [{"targets": "127.0.0.1:5000", "labels": "key=val,k1=v1"}]}`,
		},
		{
			name:    "nameEmpty",
			err:     nil,
			code:    http.StatusBadRequest,
			payload: `{"name": "", "entries": [{"targets": "127.0.0.1:5000", "labels": "key=val,k1=v1"}]}`,
		},
		{
			name:    "validationError2",
			err:     nil,
			code:    http.StatusUnprocessableEntity,
			payload: `{"name": "test", "entries": [{"targets": "127.0.0.1:5000", "labels": "keys"}]}`,
		},
		{
			name:    "notFoundError",
			err:     &db.NotFoundError{},
			code:    http.StatusNotFound,
			payload: `{"name": "test", "entries": [{"targets": "127.0.0.1:5000", "labels": "k1=v1"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/target/1", strings.NewReader(tt.payload))
			storage.returnError = tt.err
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.code, w.Code)
		})
	}
}

func Test_removeTargetHandler(t *testing.T) {
	var (
		err error
	)
	dbService, err = db.New("testdb", "")
	assert.NoError(t, err)

	router := setupRouter()

	tests := []struct {
		name string
		err  error
		code int
	}{
		{
			name: "error500",
			err:  &db.StorageError{},
			code: http.StatusInternalServerError,
		},
		{
			name: "NoErrors",
			err:  nil,
			code: http.StatusOK,
		},
		{
			name: "validationError",
			err:  &db.ValidationError{},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "notFoundError",
			err:  &db.NotFoundError{},
			code: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/api/target/1", nil)
			storage.returnError = tt.err
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.code, w.Code)
		})
	}
}

func Test_healthHandler(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_prometheusHandler(t *testing.T) {
	var (
		err error
	)
	dbService, err = db.New("testdb", "")
	assert.NoError(t, err)

	router := setupRouter()

	tests := []struct {
		name string
		err  error
		code int
	}{
		{
			name: "NoErrors",
			err:  nil,
			code: http.StatusOK,
		},
		{
			name: "validationError",
			err:  &db.ValidationError{},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "notFoundError",
			err:  &db.NotFoundError{},
			code: http.StatusNotFound,
		},
		{
			name: "error500",
			err:  &db.StorageError{},
			code: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/prom-target/1", nil)
			storage.returnError = tt.err
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.code, w.Code)
		})
	}
}
