package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

type testFn func(t *testing.T, serverURL string)

func runWithMockServer(t *testing.T, name string, fn testFn) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		testServer := newMockServer()
		defer testServer.Close()
		fn(t, testServer.URL)
	})
}

type thing struct {
	ID    int    `json:"id,omitempty"`
	Field string `json:"field"`
}

type mockServer struct {
	cache map[int]thing
}

func newMockServer() *httptest.Server {
	return httptest.NewServer(
		&mockServer{make(map[int]thing)},
	)
}

func (s *mockServer) route(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mockAuthMiddleware(s.getThing)(w, r)
	case http.MethodPost:
		mockAuthMiddleware(s.createThing)(w, r)
	case http.MethodPut:
		mockAuthMiddleware(s.uploadThing)(w, r)
	case http.MethodDelete:
		mockAuthMiddleware(s.deleteThing)(w, r)
	default:
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}(w, r)
	}
}

func (s *mockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.route(w, r)
}

func (s *mockServer) createThing(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to read body: %s", err)
		return
	}

	var t thing
	if err := json.Unmarshal(b, &t); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to decode body: %s", err)
		return
	}

	if t.Field == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "required field 'field' is unset")
		return
	}

	var mutex sync.RWMutex
	mutex.Lock()
	t.ID = len(s.cache) + 1
	s.cache[t.ID] = t
	mutex.Unlock()

	resp, err := json.Marshal(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to serialize response: %s", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

func (s *mockServer) uploadThing(w http.ResponseWriter, r *http.Request) {
	const maxMultiPartMem = 2 << 20

	if err := r.ParseMultipartForm(maxMultiPartMem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to parse multi-part form: %s", err)
		return
	}

	files := r.MultipartForm.File["file"]
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

func (s *mockServer) getThing(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["id"]
	if len(ids) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "id not provided")
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
	thing, exists := s.cache[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "thing with id %d not found", id)
		return
	}
	mutex.Unlock()

	resp, err := json.Marshal(thing)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to serialize response: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (s *mockServer) deleteThing(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["id"]
	if len(ids) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "id not provided")
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
	thing, exists := s.cache[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "thing with id %d not found", id)
		return
	}
	delete(s.cache, thing.ID)
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "successfully deleted thing")
}

// for testing setting of headers
func mockAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("auth")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "must be logged in to perform this")
			return
		}
		next(w, r)
	}
}
