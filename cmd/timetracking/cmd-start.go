package main

import (
	"flag"
	"fmt"
	"strings"
)

func commandStart(c *Command) (int, error) {
	flag.Parse()
	input := strings.Join(flag.Args(), " ")

	_, config, err := getConfig(c.l)
	if err != nil {
		return 1, err
	}

	if config == nil {
		return 1, nil
	}

	t, err := New(c.l, config)
	if err != nil {
		return 1, err
	}

	if err := t.SetUID(0); err != nil {
		return 1, err
	}

	r := config.Tasks.FuzzyFind(input, 1, true)
	if len(r) == 0 {
		return 1, fmt.Errorf("Nothing found")
	}

	task := r[0]
	entry, err := t.StartTracker(task.ProjectID, task.TaskID)
	if err != nil {
		return 0, err
	}
	c.l.Printf("Created %d", entry.ID)

	return 0, nil
}
