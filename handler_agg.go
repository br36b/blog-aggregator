package main

import (
	"context"
	"fmt"

	"github.com/br36b/blog-aggregator/internal/rssfeed"
)

func handleAggregateFeed(s *state, cmd command) error {
	// if len(cmd.args) != 1 {
	// 	return fmt.Errorf("Usage: %s <feed_url>", cmd.name)
	// }
	var feedUrl string
	if len(cmd.args) == 0 {
		feedUrl = "https://www.wagslane.dev/index.xml"
	} else {
		feedUrl = cmd.args[0]
	}

	feed, err := rssfeed.FetchFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("Unable to fetch feed: %w", err)
	}

	fmt.Printf("Feed: %+v", feed)

	return nil
}
