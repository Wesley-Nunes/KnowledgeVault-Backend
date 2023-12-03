package routes

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Router(pool *pgxpool.Pool) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/books", CreateBook(pool)).Methods("POST")
	r.HandleFunc("/details", CreateDetails(pool)).Methods("POST")

	return r
}
