package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeoHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "212.8.189.19", nil)
	if err != nil {
		t.Fatalf("count not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	h := http.HandlerFunc(geoHandler)

	h.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

}
