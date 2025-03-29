package api

import (
	"net/http"
	"context"
	"github.com/gorilla/mux"
	"strconv"
	"log"
	"fmt"
	"encoding/json"
	"github.com/shishircipher/commento/db"
)

func HandleGetComments(w http.ResponseWriter, r *http.Request) {
   //     ctx := r.Context()
        vars := mux.Vars(r)

        postID, err := strconv.Atoi(vars["postID"])
        if err != nil {
                http.Error(w, "Invalid post ID", http.StatusBadRequest)
                return
        }
	ctx := context.Background()
        comments, err := db.GetComments(ctx, postID)
        if err != nil {
                http.Error(w, fmt.Sprintf("Error retrieving comments: %v", err), http.StatusInternalServerError)
                return
        }

        // Default pagination values
//        limit := 10
//     offset := 0

        // Parse query parameters if provided
//        if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
  //              limit = l
    //    }
      //  if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil {
        //        offset = o
       // }

      // comments, err := db.GetComments(ctx, postID, limit, offset)
 //	comments := "hello, i am comments"
//	log.Println(ctx)
	log.Println(postID)
        //if err != nil {
          //      http.Error(w, fmt.Sprintf("Failed to fetch comments: %v", err), http.StatusInternalServerError)
            //    return
       // }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(comments)
}

