package main

// import (
// 	// "encoding/json"
// 	"log"
// 	"net/http"
// 	// "time"

// 	// badger "github.com/dgraph-io/badger/v4"

// 	// "mini-feed/internal/events"
// 	"mini-feed/internal/models"
// 	badgerdb "mini-feed/internal/storage/badger"
// 	"mini-feed/internal/services/fanout"
// )

import (
    "log"
    "net/http"

    "github.com/dgraph-io/badger/v4"

    "mini-feed/internal/services/fanout"
    "mini-feed/internal/services/feed"
    badgerdb "mini-feed/internal/storage/badger"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("./data/feed"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	feedStore := badgerdb.NewFeedStore(db)

	followStore := badgerdb.NewFollowStore(db)

	fanout.StartWorker(followStore, feedStore)

	handler := feed.NewHandler(feedStore)

	mux := http.NewServeMux()
	mux.HandleFunc("/feed", handler.GetFeed)

	log.Println("Feed service running on :8083")
	log.Fatal(http.ListenAndServe(":8083", mux))
}
// func GenerateID() string {
// 	return time.Now().Format("20060102150405.000000")
// }

// func main() {
//     db, _ := badger.Open(badger.DefaultOptions("./data"))
//     defer db.Close()

//     followStore := badgerdb.NewFollowStore(db)
//     postStore := badgerdb.NewPostStore(db)
//     feedStore := badgerdb.NewFeedStore(db)

//     // fanout worker
//     fanout.StartWorker(followStore, feedStore)

//     mux := http.NewServeMux()

//     // post creation
//     mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
//         user := r.URL.Query().Get("user")
//         text := r.URL.Query().Get("text")
//         post := &models.Post{ID: GenerateID(), UserID: user, Text: text}
//         postStore.CreatePost(post)
//         events.PostEventChannel <- post
//         w.Write([]byte("ok"))
//     })

//     // feed fetch
//     mux.HandleFunc("/feed", func(w http.ResponseWriter, r *http.Request) {
//         user := r.URL.Query().Get("user")
//         posts, _ := feedStore.GetFeed(user, 20)
//         json.NewEncoder(w).Encode(map[string]any{"user": user, "posts": posts})
//     })

//     log.Println("Server running on :8080")
//     log.Fatal(http.ListenAndServe(":8080", mux))
// }