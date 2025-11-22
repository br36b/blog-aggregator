package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/br36b/blog-aggregator/internal/database"
)

func handleBrowseRssPosts(s *state, cmd command, user database.User) error {
	itemLimitBase32 := int32(2)

	if len(cmd.args) > 1 {
		return fmt.Errorf("Usage: %s [post_limit]", cmd.name)
	}

	if len(cmd.args) == 1 {
		itemLimitBase64, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("Invalid post limit set: %w", err)
		}

		itemLimitBase32 = int32(itemLimitBase64)
	}

	getPostsParams := database.GetPostsParams{
		UserID: user.ID,
		Limit:  itemLimitBase32,
	}

	posts, err := s.db.GetPosts(context.Background(), getPostsParams)
	if err != nil {
		return fmt.Errorf("Failed to get posts for user: %w", err)
	}

	fmt.Printf("%d latest post(s):\n", itemLimitBase32)

	for _, post := range posts {
		fmt.Println("\nTitle:", post.Title)
		fmt.Println("Description:", post.Description)
		fmt.Println("Link:", post.Url)
		fmt.Println("Publication Date:", post.PublishedAt)
		fmt.Println()
	}

	return nil
}
