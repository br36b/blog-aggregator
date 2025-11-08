package main

import (
	"context"
	"fmt"

	"github.com/br36b/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Failed to get user entry: %w", err)
		}

		return handler(s, c, user)
	}
}
