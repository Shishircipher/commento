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

	// Initialize database
	database, err := db.InitDb()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer database.Close() // Ensure the connection pool is closed when main exits
	// Create API router and inject database connection

	r := mux.NewRouter()
	r.HandleFunc("/comments/{postID}", func(w http.ResponseWriter, r *http.Request) {
		api.HandleGetComments(w, r, database) // Inject DB into handler
	}).Methods("GET")
	r.HandleFunc("/comments/{postID}/{commentID}", func(w http.ResponseWriter, r *http.Request) {
                api.HandleGetComments(w, r, database) // Inject DB into handler
        }).Methods("GET")

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


