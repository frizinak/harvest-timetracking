package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/frizinak/harvest-timetracking/config"
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

func getConfig(l *log.Logger) (*config.ConfigLoader, *Config, error) {
	confLoader, err := config.DotFile(
		".timetracking",
		&Config{
			"-- your account id --",
			"-- your forecast account id (optional)--",
			defaultToken,
			[]string{},
			[]string{"saturday", "sunday"},
			nil,
			nil,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	conf := &Config{}
	if err := confLoader.Read(conf); err != nil {
		if os.IsNotExist(err) {
			l.Printf(
				"Config file %s does not exist, creating example. [https://id.getharvest.com/developers to create an access token]",
				confLoader.Path(),
			)
			return nil, nil, confLoader.CreateDefault()
		}

		l.Printf("Failed to parse config")
		return nil, nil, err
	}

	if conf.Token == defaultToken {
		l.Printf(
			"You should fill in your access token and account id in '%s'",
			confLoader.Path(),
		)

		return nil, nil, nil
	}

	return confLoader, conf, nil
}

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

	exit, err := c.Run(arg)
	if err != nil {
		l.Println(err)
	}

	os.Exit(exit)
}
