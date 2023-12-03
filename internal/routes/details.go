package routes

import (
	"database/sql/driver"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReadingStatus string

const (
	Read     ReadingStatus = "Read"
	Reading  ReadingStatus = "Reading"
	Wishlist ReadingStatus = "Wishlist"
)

type Details struct {
	Status         ReadingStatus `json:"status"`
	Pages          int           `json:"pages"`
	CurrentPage    int           `json:"currentPage"`
	PercentRead    int           `json:"percentRead"`
	StartReadingAt time.Time     `json:"startReadingAt"`
	EndReadingAt   time.Time     `json:"endReadingAt"`
}

func (rs *ReadingStatus) Scan(value interface{}) error {
	*rs = ReadingStatus(value.(string))
	return nil
}

func (rs ReadingStatus) Value() (driver.Value, error) {
	return string(rs), nil
}

func CreateDetails(dbConn *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlInsert := `INSERT INTO details (
			status,
			pages,
			current_page,
			percent_read,
			start_reading_at,
			end_reading_at
		)
		VALUES ($1, $2, $3, $4, $5, $6);`
		var details Details

		if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
			http.Error(w, "Error: Invalid request body", http.StatusBadRequest)
			return
		}

		if details.Status == "" || details.Pages == 0 {
			http.Error(w, "Error: status, and pages are required fields", http.StatusUnprocessableEntity)
			return
		}

		_, err := dbConn.Exec(r.Context(), sqlInsert, details.Status, details.Pages, details.CurrentPage, details.PercentRead, details.StartReadingAt, details.EndReadingAt)
		if err != nil {
			if strings.Contains(err.Error(), "SQLSTATE 22P02") {
				http.Error(w, "Error: The possible values for status are 'Read' or 'Reading' or 'Wishlist'.", http.StatusUnprocessableEntity)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Printf("Create successfully: details\n")
	}
}
