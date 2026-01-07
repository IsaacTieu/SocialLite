package main

import (
	"log"
	"net/http"
	badgerdb "mini-feed/internal/storage/badger"
	"mini-feed/internal/services/post"
)

func main() {
	db, err := badgerdb.Open("./data/post")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := badgerdb.NewPostStore(db)
	handler := post.NewHandler(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/posts", handler.CreatePost)

	log.Println("Post service listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", mux))
}