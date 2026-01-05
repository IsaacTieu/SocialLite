package main

import (
	"log"
	"net/http"

	"mini-feed/internal/services/user"
	badgerdb "mini-feed/internal/storage/badger"
)

func main() {
	db, err := badgerdb.Open("./data")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := badgerdb.NewUserStore(db)
	handler := user.NewHandler(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func (w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/users", handler.CreateUser)
	mux.HandleFunc("/follow", handler.Follow)

	log.Println("user service listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}