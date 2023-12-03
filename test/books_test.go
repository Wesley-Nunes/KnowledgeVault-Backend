package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wesley-Nunes/KnowledgeVault-Backend/internal/durable"
	"github.com/Wesley-Nunes/KnowledgeVault-Backend/internal/routes"
)

func TestCreateBook_InvalidPayload(t *testing.T) {
	rr := httptest.NewRecorder()
	pool := durable.GetConnPool()
	defer pool.Close()
	handler := http.HandlerFunc(routes.CreateBook(pool))

	req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
}
func TestCreateBook_BookAlreadyRegistered(t *testing.T) {
	rr := httptest.NewRecorder()
	pool := durable.GetConnPool()
	defer pool.Close()
	handler := http.HandlerFunc(routes.CreateBook(pool))

	book := []byte(`{
		"title": "Frankenstein",
		"author": "Mary Shelley",
		"pages": 402
	}`)
	req, err := http.NewRequest("POST", "/books", bytes.NewBuffer(book))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}
}

func TestCreateBook_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	pool := durable.GetConnPool()
	defer pool.Close()
	handler := http.HandlerFunc(routes.CreateBook(pool))

	book := []byte(`{
		"title": "Title fake to test",
		"author": "Author fake to test",
		"pages": 999999999
	}`)
	req, err := http.NewRequest("POST", "/books", bytes.NewBuffer(book))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}
