package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.name)
	}

	err := s.db.ResetUserDb(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset users database: %w", err)
	}

	fmt.Println("Successfully reset the users database")

	return nil
}
