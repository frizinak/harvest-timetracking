package main

import (
	"log"
	"sort"
)

type Cmd struct {
	Description string
	Command     func(c *Command) (int, error)
}

type Command struct {
	l        *log.Logger
	commands map[string]*Cmd
}

func (c *Command) Run(arg string) (int, error) {
	cmd := c.commands[arg]
	if cmd == nil {
		return commandHelp(c)
	}

	return cmd.Command(c)
}

func commandHelp(c *Command) (int, error) {
	cmds := make([]string, 0, len(c.commands))
	for i := range c.commands {
		cmds = append(cmds, i)
	}

	sort.SliceStable(
		cmds,
		func(i, j int) bool {
			return cmds[i] < cmds[j]
		},
	)

	c.l.Println("Available commands:")
	for _, cmd := range cmds {
		c.l.Printf("  %-20s - %s\n", cmd, c.commands[cmd].Description)
	}

	return 0, nil
}
