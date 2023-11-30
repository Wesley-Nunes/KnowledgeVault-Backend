package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func shutdownServer(srv *http.Server, wait time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server shutdown error:", err)
	}

	log.Println("Shutting down")
	os.Exit(0)
}

func GetConnPool() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(os.Getenv("KNOWLEDGEVAULT_DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse database URL: %v\n", err)
		os.Exit(1)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	log.Println("Connected")

	return pool
}

func main() {
	// Load the server
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully waits for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	pool := GetConnPool()
	defer pool.Close()
	r := mux.NewRouter()
	r.HandleFunc("/books", CreateBook(pool)).Methods("POST")

	srv := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	shutdownServer(srv, wait)
}
