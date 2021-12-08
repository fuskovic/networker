package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
)

// MockObject is an arbitrary object that can be encoded/decoded
// and it's primary function is to provide a shape for the mock
// web server to use for its transport and store layers.
type MockObject struct {
	ID    int    `json:"id,omitempty"`
	Field string `json:"field"`
}

type mockServer struct {
	cache map[int]MockObject
}

// newMockServer returns test server with pre-configured endpoints
// that uses an in-memory cache for its store layer.
func newMockServer() *httptest.Server {
	return httptest.NewServer(
		&mockServer{make(map[int]MockObject)},
	)
}

func (s *mockServer) route(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mockAuthMiddleware(s.getObject)(w, r)
	case http.MethodPost:
		mockAuthMiddleware(s.createObject)(w, r)
	case http.MethodPut:
		mockAuthMiddleware(s.uploadFile)(w, r)
	case http.MethodDelete:
		mockAuthMiddleware(s.deleteObject)(w, r)
	default:
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}(w, r)
	}
}

// ServeHTTP routes handles routing all requests to their respective handlers and implements http.Handler.
func (s *mockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.route(w, r)
}

func (s *mockServer) createObject(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to read body: %s", err)
		return
	}

	var o MockObject
	if err := json.Unmarshal(b, &o); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to decode body: %s", err)
		return
	}

	if o.Field == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "required field 'field' is unset")
		return
	}

	var mutex sync.RWMutex
	mutex.Lock()
	o.ID = len(s.cache) + 1
	s.cache[o.ID] = o
	mutex.Unlock()

	resp, err := json.Marshal(o)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to serialize response: %s", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

func (s *mockServer) uploadFile(w http.ResponseWriter, r *http.Request) {
	const maxMultiPartMem = 2 << 20

	if err := r.ParseMultipartForm(maxMultiPartMem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to parse multi-part form: %s", err)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "failed to retrieve multi-part form files")
		return
	}

	if len(files) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "expected 2 files")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "uploaded files")
}

func (s *mockServer) getObject(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["id"]
	if len(ids) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "id not provided for get")
		return
	}

	id, err := strconv.Atoi(ids[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to convert id string to int: %s", err)
		return
	}

	var mutex sync.RWMutex
	mutex.Lock()
	o, exists := s.cache[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "object with id %d not found", id)
		return
	}
	mutex.Unlock()

	resp, err := json.Marshal(o)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to serialize response: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (s *mockServer) deleteObject(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["id"]
	if len(ids) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "id not provided for delete")
		return
	}

	id, err := strconv.Atoi(ids[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to convert id string to int: %s", err)
		return
	}

	var mutex sync.RWMutex
	mutex.Lock()
	o, exists := s.cache[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "object with id %d not found", id)
		return
	}
	delete(s.cache, o.ID)
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "successfully deleted thing")
}

// for testing setting of headers
func mockAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "must be logged in to perform this")
			return
		}
		if auth != "Bearer doesntmatter" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
		next(w, r)
	}
}
