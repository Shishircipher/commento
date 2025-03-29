package main

import (
	"fmt"
	"log"
//	"database/sql"
//	"github.com/jackc/pgx"
	"net/http"
	"github.com/shishircipher/commento/db"
	"github.com/shishircipher/commento/api"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
// Initialize the database connection
	db.InitDb()
	defer db.CloseDb() // Correctly closing the database

	r := mux.NewRouter()
	r.HandleFunc("/comments/{postID}", api.HandleGetComments).Methods("GET")

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

