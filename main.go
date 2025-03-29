package main

import (
	"fmt"
	"log"
//	"database/sql"
//	"github.com/jackc/pgx"
	"github.com/shishircipher/commento/db"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db.InitDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/comments/{postID}", api.handleGetComments).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}

