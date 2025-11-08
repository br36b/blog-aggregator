package main

import (
	"context"
	"fmt"
	"time"

	"github.com/br36b/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handleFollowRssFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.name)
	}

	argUrl := cmd.args[0]

	// Get Feed
	feedEntry, err := s.db.GetFeedByUrl(context.Background(), argUrl)
	if err != nil {
		return fmt.Errorf("No feed matches the URL (%s): %w", argUrl, err)
	}

	// Follow Feed
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedEntry.ID,
	}
	s.db.CreateFeedFollow(context.Background(), feedFollowParams)

	return nil
}

func handleGetFollowedRssFeeds(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.name)
	}

	// Get Followed Feeds
	followedFeedsEntry, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Unable to get followed feeds for user (%s): %w", user.Name, err)
	}

	fmt.Println("Feeds Following:")

	if len(followedFeedsEntry) == 0 {
		fmt.Println("\tNo feeds are currently being followed by this user.")
	}

	for _, feed := range followedFeedsEntry {
		fmt.Printf("\tFeed: %s\n", feed.FeedName)
	}

	return nil
}
