package badgerdb

import (
	"github.com/dgraph-io/badger/v4"
	"encoding/json"
)

type FeedStore struct {
	db *badger.DB
}

func NewFeedStore(db *badger.DB) *FeedStore {
	return &FeedStore{db: db}
}

func (s *FeedStore) AddToFeed(userID, postID string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		key := []byte("feed:" + userID)
		var posts []string

		item, err := txn.Get(key)

		if err == nil {
			val, _ := item.ValueCopy(nil)
			_ = json.Unmarshal(val, &posts)
		}

		posts = append([]string{postID}, posts...) // pre-pend the new post

		data, _ := json.Marshal(posts)
		return txn.Set(key, data)
	})
}

func (s *FeedStore) GetFeed(userID string, limit int) ([]string, error) {
	var posts []string

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("feed:" + userID))
		if err != nil {
			return nil
		}

		val, _ := item.ValueCopy(nil)
		return json.Unmarshal(val, &posts)
	})

	if limit > 0 && len(posts) > limit {
		posts = posts[:limit]
	}
	return posts, err
}