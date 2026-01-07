package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"mini-feed/internal/middleware"
	"github.com/redis/go-redis/v9"
	"mini-feed/internal/services/user"
	badgerdb "mini-feed/internal/storage/badger"
)

func main() {
	db, err := badgerdb.Open("./data/user")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	store := badgerdb.NewUserStore(db)
	handler := user.NewHandler(store, rdb)

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/users", handler.CreateUser)
	mux.HandleFunc("/follow", handler.Follow)

	log.Println("user service listening on :8081")

	wrappedMux := middleware.LoggerAndMetrics(mux)
	log.Fatal(http.ListenAndServe(":8081", wrappedMux))
}
