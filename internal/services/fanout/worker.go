package fanout

import (
	"log"
	"mini-feed/internal/events"
)

type FollowStore interface {
	GetFollowers(usedID string) []string
}

type FeedStore interface {
	AddToFeed(userID, postID string) error
}

func StartWorker(followStore FollowStore, feedStore FeedStore) {
	go func() {
		log.Println("Fanout worker started")

		for post := range events.PostEventChannel {
			followers := followStore.GetFollowers(post.UserID)

			for _, follower := range followers {
				err := feedStore.AddToFeed(follower, post.ID)
				if err != nil {
					log.Println("Fanout error:", err)
				}
			}
		}
	} ()
}