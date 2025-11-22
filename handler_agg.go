package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/br36b/blog-aggregator/internal/database"
	"github.com/br36b/blog-aggregator/internal/rssfeed"
	"github.com/google/uuid"
)

func generatePostParams(item rssfeed.RSSItem, parentFeed database.Feed) (database.CreatePostParams, error) {
	// fmt.Printf("\n\tPost: %s - %s\n\n", item.Title, item.PubDate)
	parsedPublishDate, err := time.Parse(time.RFC1123, item.PubDate)
	if err != nil {
		return database.CreatePostParams{}, err
	}

	nullableDescription := sql.NullString{
		String: item.Description,
		Valid:  item.Description != "",
	}

	createPostParams := database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       item.Title,
		Url:         item.Link,
		Description: nullableDescription,
		PublishedAt: parsedPublishDate,
		FeedID:      parentFeed.ID,
	}

	return createPostParams, nil
}

func scrapeFeeds(s *state) {
	log.Println("\n\nScraping feeds:")
	// Get next feed to fetch
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Unable to find feed to fetch: ", err)
		return
	}

	// Mark feed as fetched
	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		log.Println("Unable to mark feed as fetched: ", err)
		return
	}

	// Fetch feed RSS data by URL
	rssFeed, err := rssfeed.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		log.Println("Unable to poll feed data: ", err)
		return
	}

	// Save feed items as posts
	log.Printf("\n\nUpdated feed: %s - %s\n", rssFeed.Channel.Title, rssFeed.Channel.Description)
	for _, item := range rssFeed.Channel.Item {
		postParams, err := generatePostParams(item, nextFeed)
		if err != nil {
			log.Println("Unable to create post parameters: ", err)
			continue
		}

		err = s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			log.Println("Unable to create post: ", err)
		}
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
