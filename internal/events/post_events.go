package events

import "mini-feed/internal/models"

var PostEventChannel = make(chan *models.Post, 100)