package main

import (
	"context"
	"fmt"
	"time"

	"github.com/br36b/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Username must be provided for this command")
	}

	if len(cmd.args) > 1 {
		return fmt.Errorf("Too many arguments provided for this command")
	}

	username := cmd.args[0]

	newUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	userEntry, err := s.db.CreateUser(context.Background(), newUserParams)
	if err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to save user: %w", err)
	}

	fmt.Printf("Successfully created user: %+v\n", userEntry)

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Username must be provided for this command")
	}

	if len(cmd.args) > 1 {
		return fmt.Errorf("Too many arguments provided for this command")
	}

	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User '%s' was not found: %w", username, err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to login as '%s': %w", username, err)
	}

	fmt.Printf("Successfully logged in as: %s\n", username)

	return nil
}

func handleGetUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.name)
	}

	userEntries, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset users database: %w", err)
	}

	fmt.Println("List of all users:")
	for _, user := range userEntries {
		lineOutput := fmt.Sprintf("  * %s", user.Name)
		if user.Name == s.cfg.CurrentUserName {
			lineOutput += " (current)"
		}

		fmt.Println(lineOutput)
	}

	fmt.Println("Successfully reset the users database")

	return nil
}
