package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wesley-Nunes/KnowledgeVault-Backend/internal/durable"
	"github.com/Wesley-Nunes/KnowledgeVault-Backend/internal/routes"
)

func TestCreateDetails_invalidPayload(t *testing.T) {
	rr := httptest.NewRecorder()
	pool := durable.GetConnPool()
	defer pool.Close()
	handler := http.HandlerFunc(routes.CreateDetails(pool))

	req, err := http.NewRequest("POST", "/details", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
}

func TestCreateDetails_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	pool := durable.GetConnPool()
	defer pool.Close()
	handler := http.HandlerFunc(routes.CreateDetails(pool))

	details := []byte(`{
		"status": "Wishlist",
		"pages": 999999999
	}`)
	req, err := http.NewRequest("POST", "/details", bytes.NewBuffer(details))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}
