package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frioux/leatherman/internal/middleware"
	"github.com/frioux/leatherman/internal/testutil"
)

func TestLog(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	var inner http.HandlerFunc = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(404)
	}
	rr := httptest.NewRecorder()
	handler := middleware.Adapt(inner, middleware.Log(buf))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	d := json.NewDecoder(buf)
	var x struct{ StatusCode int }
	if err = d.Decode(&x); err != nil {
		panic(err)
	}

	testutil.Equal(t, x.StatusCode, http.StatusNotFound, "status code recorded")
}
