package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var v = "unknown"

const (
	endOfWeek    = "end-of-week"
	nextWeek     = "next-week"
	defaultToken = "-- your account token --"

	groupByDay   = "day"
	groupByWeek  = "week"
	groupByMonth = "month"
	groupByYear  = "year"
)

type Duration time.Duration

func (d Duration) String() string {
	s := time.Duration(d)
	h := s / time.Hour
	m := (s % time.Hour) / time.Minute
	return fmt.Sprintf("%dh%02d", h, m)
}

func main() {
	arg := ""
	if len(os.Args) > 1 && os.Args[1][0] != '-' {
		arg = os.Args[1]
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	l := log.New(os.Stdout, "", 0)
	c := &Command{
		l:        l,
		commands: make(map[string]*Cmd),
	}
	c.commands["version"] = &Cmd{"print version", commandVersion}
	c.commands["help"] = &Cmd{"print list of commands", commandHelp}
	c.commands["tracking"] = &Cmd{"show tracked hours", commandTracking}
	c.commands["off"] = &Cmd{"get a list of days off using the forecast api", commandDaysOff}
	c.commands["tasks"] = &Cmd{"get a list of projects and their tasks", commandTasks}
	c.commands["start"] = &Cmd{"Start a timetracker", commandStart}

	exit, err := c.Run(arg)
	if err != nil {
		l.Println(err)
	}

	os.Exit(exit)
}
