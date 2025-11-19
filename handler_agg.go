package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/br36b/blog-aggregator/internal/rssfeed"
)

func scrapeFeeds(s *state) {
	fmt.Println("\n\n\nScraping feeds at interval:", time.Now())
	// Get next feed to fetch
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Unable to find feed to fetch: %w", err)
		return
	}

	// Mark feed as fetched
	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		log.Println("Unable to mark feed as fetched: %w", err)
		return
	}

	// Fetch feed RSS data by URL
	rssFeed, err := rssfeed.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		log.Println("Unable to poll feed data: %w", err)
		return
	}

	// Iterate over the feed items
	fmt.Printf("\nUpdated feed: %s - %s\n", rssFeed.Channel.Title, rssFeed.Channel.Description)
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("\tPost: %s - %s\n\n", item.Title, item.PubDate)
	}
}

func handleAggregateFeed(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <time_between_requests>", cmd.name)
	}

	// Parse interval
	intervalArg := cmd.args[0]
	interval, err := time.ParseDuration(intervalArg)
	if err != nil {
		return fmt.Errorf("Unable to parse interval: %w", err)
	}

	// Poll for feed updates
	ticker := time.NewTicker(interval)

	fmt.Println("Polling feeds every:", interval)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
