package main

import (
	"fmt"
	"os"

	"github.com/br36b/blog-aggregator/internal/config"
)

type state struct {
	cfg config.Config
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Username must be provided for this command")
	}

	if len(cmd.args) > 1 {
		return fmt.Errorf("Too many arguments provided for this command")
	}

	username := cmd.args[0]

	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to login and save user: %v", err)
	}

	fmt.Printf("Successfully logged in as: %s\n", username)

	return nil
}

func main() {
	// Basic app initialization
	dbConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	appState := &state{
		cfg: dbConfig,
	}

	appCommands := &commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	appCommands.register("login", handlerLogin)

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
