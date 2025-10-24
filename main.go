package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/br36b/blog-aggregator/internal/config"
	"github.com/br36b/blog-aggregator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if callback, ok := c.commandMap[cmd.name]; ok {
		return callback(s, cmd)
	} else {
		return fmt.Errorf("No command found with given name: %s", cmd.name)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

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
		return fmt.Errorf("Failed to create user: %v", err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to save user: %v", err)
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
		return fmt.Errorf("User '%s' was not found: %v", username, err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to login as '%s': %v", username, err)
	}

	fmt.Printf("Successfully logged in as: %s\n", username)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.name)
	}

	err := s.db.ResetUserDb(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset users database: %v", err)
	}

	fmt.Println("Successfully reset the users database")

	return nil
}

func main() {
	// Basic app initialization
	dbConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("postgres", dbConfig.DbUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	appState := &state{
		cfg: &dbConfig,
		db:  dbQueries,
	}

	appCommands := &commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)
	appCommands.register("reset", handlerReset)

	// Command processing
	commandArgs := os.Args

	if len(commandArgs) < 2 {
		fmt.Fprintln(os.Stderr, "Error: At least two arguments are required")
		os.Exit(1)
	}

	userCommand := command{
		name: commandArgs[1],
		args: commandArgs[2:],
	}

	err = appCommands.run(appState, userCommand)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	// fmt.Printf("%+v\n\n%+v\n", appState, commandArgs)
}
