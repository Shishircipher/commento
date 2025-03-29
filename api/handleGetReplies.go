package api

import (
	"net/http"
//	"context"
	"github.com/gorilla/mux"
	"strconv"
//	"log"
	"fmt"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shishircipher/commento/db"
)

func handleGetReplies(w http.ResponseWriter, r *http.Request, database *pgxpool.Pool) {
        ctx := r.Context()
        vars := mux.Vars(r)

        postID, err1 := strconv.Atoi(vars["postID"])
        commentID, err2 := strconv.Atoi(vars["commentID"])
        if err1 != nil || err2 != nil {
                http.Error(w, "Invalid post ID or comment ID", http.StatusBadRequest)
                return
        }

        // Default depth limit and pagination
        depthLimit := 3
        limit := 10
        offset := 0

        // Parse query parameters if provided
        if d, err := strconv.Atoi(r.URL.Query().Get("depth")); err == nil {
                depthLimit = d
        }
        if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
                limit = l
        }
        if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil {
                offset = o
        }

        replies, err := db.GetReplies(ctx, database,postID, commentID, depthLimit, limit, offset)
        if err != nil {
                http.Error(w, fmt.Sprintf("Failed to fetch replies: %v", err), http.StatusInternalServerError)
                return
        }
	w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(replies)
}
