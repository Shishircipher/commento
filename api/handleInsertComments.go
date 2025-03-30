package api

import (
	"net/http"
//	"context"
	//"github.com/gorilla/mux"
//	"strconv"
//	"log"
	"fmt"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shishircipher/commento/db"
)
func HandleInsertComment(w http.ResponseWriter, r *http.Request, database *pgxpool.Pool) {
        ctx := r.Context()

        var newComment db.NewComment
        if err := json.NewDecoder(r.Body).Decode(&newComment); err != nil {
                http.Error(w, "Invalid request payload", http.StatusBadRequest)
                return
        }

        commentID, err := db.InsertComment(ctx, database, newComment)

        if err != nil {
                http.Error(w, fmt.Sprintf("Failed to insert comment: %v", err), http.StatusInternalServerError)
                return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]int{"comment_id": commentID})
}
