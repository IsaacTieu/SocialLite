package badgerdb

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
)

type FollowStore struct {
	db *badger.DB
}

func NewFollowStore(db *badger.DB) *FollowStore {
	return &FollowStore{db: db}
}

// follower follows user
func (s *FollowStore) AddFollow(follower, user string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		key := []byte("followers:" + user)

		var followers []string
		item, err := txn.Get(key)
		if err == nil {
			val, _ := item.ValueCopy(nil)
			_ = json.Unmarshal(val, &followers)
		}

		followers = append(followers, follower)

		data, _ := json.Marshal(followers)
		return txn.Set(key, data)
	})
}

func (s *FollowStore) GetFollowers(user string) []string {
	var followers []string

	_ = s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("followers:" + user))
		if err != nil {
			return nil
		}
		val, _ := item.ValueCopy(nil)
		_ = json.Unmarshal(val, &followers)
		return nil
	})

	return followers
}
