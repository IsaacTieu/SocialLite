package fanout

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
	"mini-feed/internal/models"
)

var fanoutOps = promauto.NewCounter(prometheus.CounterOpts{
	Name: "fanout_processed_ops_total",
	Help: "The total number of processed fanout operations",
})

type FollowStore interface {
	GetFollowers(usedID string) []string
	AddFollow(follower, user string) error
}

type FeedStore interface {
	AddToFeed(userID, postID string) error
}

func StartWorker(followStore FollowStore, feedStore FeedStore, rdb *redis.Client) {
	ctx := context.Background()

	sub := rdb.Subscribe(ctx, "post_created", "user_followed")

	go func() {
		// log.Println("Redis fanout worker started")
		slog.Info("Redis fanout worker started")
		ch := sub.Channel()

		for msg := range ch {
			switch msg.Channel {
			case "user_followed":
				var event map[string]string
				if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
					slog.Error("error unmarshaling event", "error", err)
					continue
				}
				if err := followStore.AddFollow(event["follower"], event["followee"]); err != nil {
					slog.Error("Error syncing follow:", "error", err)
				} else {
					slog.Info("synced follow", 
						"follower", event["follower"], 
						"followee", event["followee"],
					)
				}

			case "post_created":
				var post models.Post
				if err := json.Unmarshal([]byte(msg.Payload), &post); err != nil {
					slog.Error("error unmarshaling post event:", "error", err)
					continue
				}

				followers := followStore.GetFollowers(post.UserID)
				for _, follower := range followers {
					err := feedStore.AddToFeed(follower, post.ID)
					if err != nil {
						slog.Error("failed to fanout post", 
							"post_id", post.ID, 
							"follower", follower, 
							"error", err,
						)
					}
				}
				fanoutOps.Inc()
				slog.Info(
					"fanned out post", 
					"post_id", post.ID, 
					"user_id", post.UserID, 
					"follower_count", len(followers),
				)
			}
		}
	} ()
}