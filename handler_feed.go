package main

import (
	"context"
	"fmt"
	"time"

	"github.com/br36b/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handleAddRssFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Usage: %s <feed_title> <feed_url>", cmd.name)
	}

	// Get Current User
	username := s.cfg.CurrentUserName
	userEntry, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User '%s' was not found: %v", username, err)
	}

	// Save Feed
	feedName, feedUrl := cmd.args[0], cmd.args[1]

	newFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    userEntry.ID,
	}

	feedEntry, err := s.db.CreateFeed(context.Background(), newFeedParams)
	if err != nil {
		return fmt.Errorf("Unable to create feed: %w", err)
	}

	fmt.Printf("Feed: %+v", feedEntry)

	return nil
}
