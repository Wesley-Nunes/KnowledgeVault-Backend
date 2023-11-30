package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Run tests
	exitCode := m.Run()

	// Run cleanup after all tests
	successCleanUp()

	// Exit with the result of the test run
	os.Exit(exitCode)
}

func TestCreateBook_InvalidPayload(t *testing.T) {
	rr := httptest.NewRecorder()
	pool := GetConnPool()
	defer pool.Close()
	handler := http.HandlerFunc(CreateBook(pool))

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
	pool := GetConnPool()
	defer pool.Close()
	handler := http.HandlerFunc(CreateBook(pool))

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
	pool := GetConnPool()
	defer pool.Close()
	handler := http.HandlerFunc(CreateBook(pool))

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

func successCleanUp() {
	sqlDelete := "DELETE FROM books WHERE author = 'Author fake to test' AND title = 'Title fake to test' AND pages = 999999999;"
	pool := GetConnPool()
	defer pool.Close()

	_, err := pool.Exec(context.Background(), sqlDelete)
	if err != nil {
		// Handle the error, log it, or return it as needed
		fmt.Println("Error cleaning up:", err)
	}
}
