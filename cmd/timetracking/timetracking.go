package main

import (
	"flag"
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
)

func getConfig(l *log.Logger) (*Config, error) {
	confLoader, err := config.DotFile(
		".timetracking",
		&Config{"-- your account id --", defaultToken, []string{}, nil},
	)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	if err := confLoader.Read(conf); err != nil {
		if os.IsNotExist(err) {
			l.Printf(
				"Config file %s does not exist, creating example. [https://id.getharvest.com/developers to create an access token]",
				confLoader.Path(),
			)
			return nil, confLoader.CreateDefault()
		}

		l.Printf("Failed to parse config")
		return nil, err
	}

	if conf.Token == defaultToken {
		l.Printf(
			"You should fill in your access token and account id in '%s'",
			confLoader.Path(),
		)

		return nil, nil
	}

	return conf, nil
}

type Duration time.Duration

func (d Duration) String() string {
	s := time.Duration(d)
	h := s / time.Hour
	m := (s % time.Hour) / time.Minute
	return fmt.Sprintf("%dh%02d", h, m)
}

func main() {
	l := log.New(os.Stdout, "", 0)
	version := flag.Bool("v", false, "Print version and exit")
	var userID int
	var days int
	var customCapacity int
	var customDate string
	var onlyWorkedDays bool
	flag.IntVar(&userID, "uid", 0, "The user id of the user to fetch time entries for")
	flag.IntVar(&days, "days", 20, "Amount of days to retrieve time entries for")
	flag.IntVar(&customCapacity, "hours", 0, "Amount of hours in a single workweek (default: from harvest api)")
	flag.BoolVar(&onlyWorkedDays, "worked", false, "Only track days that have tracking entries")
	flag.StringVar(
		&customDate,
		"from",
		"",
		fmt.Sprintf(
			"Custom date to start at [YYYY-MM-DD or %s or %s]",
			endOfWeek,
			nextWeek,
		),
	)
	//userName := flag.String("user", "", "The user name to fetch time entries for")
	flag.Parse()

	if *version {
		l.Println(v)
		os.Exit(0)
	}

	config, err := getConfig(l)
	if err != nil {
		l.Println(err)
		os.Exit(1)
	}

	if config == nil {
		os.Exit(1)
	}

	t, err := New(l, config)
	if err != nil {
		l.Println(err)
		os.Exit(1)
	}

	from := time.Now()

	if err := t.SetUID(userID); err != nil {
		l.Println(err)
		os.Exit(1)
	}

	capacity := Duration(t.User().Capacity())
	if customCapacity != 0 {
		capacity = Duration(customCapacity) * Duration(time.Hour)
	}
	daysCapacity := Duration(float64(capacity) * float64(days) / 5)

	switch {
	case customDate == endOfWeek || customDate == nextWeek:
		wd := from.Weekday() - 1
		if wd < 0 {
			wd = 6
		}
		from = from.AddDate(0, 0, 6-int(wd))
		if customDate == nextWeek {
			from = from.AddDate(0, 0, 7)
		}

	case customDate != "":
		f, err := time.Parse("2006-01-02", customDate)
		if err != nil {
			l.Printf("Invalid date '%s' expected YYYY-mm-dd", customDate)
			os.Exit(1)
		}
		from = f
	}

	l.Printf(
		"Running for %s %s\nID: %d\nWeek: %s\nOver %d days: %s\nFrom: %s\n\n",
		t.User().FirstName,
		t.User().LastName,
		t.User().ID,
		capacity,
		days,
		daysCapacity,
		from.Format("Mon Jan 02 2006"),
	)

	v, err := t.GetRecentDays(days, from, !onlyWorkedDays)
	if err != nil {
		panic(err)
	}

	var sum time.Duration
	for _, e := range v.Days(true) {
		l.Printf(
			"%s - %5s",
			e.SpentDate.Format("Mon Jan 02 2006"),
			Duration(e.Hours),
		)
		sum += e.Hours
	}

	diff := daysCapacity - Duration(sum)
	diffStr := fmt.Sprintf("%s remaining...", diff)
	if diff < 0 {
		diffStr = "target reached!"
	}

	l.Printf(
		"\nTotal: %s / %s (%.2f%%)\n%s",
		Duration(sum),
		daysCapacity,
		100*float64(sum)/float64(daysCapacity),
		diffStr,
	)
}
