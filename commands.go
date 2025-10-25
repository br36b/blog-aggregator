package main

import (
	"fmt"
)

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
