package main

import "errors"

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) Run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("unknown command: " + cmd.Name)
	}
	return f(s, cmd)
}

func (c *commands) Register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
