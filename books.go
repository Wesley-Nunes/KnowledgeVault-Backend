package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
	Cover  string `json:"cover"`
}

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
	Cover  string `json:"cover"`
}

func CreateBook(dbConn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlInsert := "INSERT INTO books (title, author, pages, cover) VALUES ($1, $2, $3, $4) RETURNING id;"
		var reqBook BookRequest
		var bookId int

		if err := json.NewDecoder(r.Body).Decode(&reqBook); err != nil {
			http.Error(w, "Error: Invalid request body", http.StatusBadRequest)
			return
		}

		if reqBook.Title == "" || reqBook.Author == "" || reqBook.Pages == 0 {
			http.Error(w, "Error: Title, Author, and Pages are required fields", http.StatusUnprocessableEntity)
			return
		}

		err := dbConn.QueryRow(r.Context(), sqlInsert, reqBook.Title, reqBook.Author, reqBook.Pages, reqBook.Cover).Scan(&bookId)
		if err != nil {
			if strings.Contains(err.Error(), "SQLSTATE 23505") {
				http.Error(w, "Error: Book already registered.", http.StatusConflict)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		book := Book{
			Id:     bookId,
			Title:  reqBook.Title,
			Author: reqBook.Author,
			Pages:  reqBook.Pages,
			Cover:  reqBook.Cover,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(book)
		log.Printf("Create successfully:\n%+v\n", book)
	}
}
