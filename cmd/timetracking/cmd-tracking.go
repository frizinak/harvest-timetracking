package main

import (
	"flag"
	"fmt"
	"time"
)

func commandTracking(c *Command) (int, error) {
	var userID int
	var days int
	var customCapacity int
	var customDate string
	var onlyWorkedDays bool
	var group string
	flag.IntVar(&userID, "uid", 0, "The user id of the user to fetch time entries for")
	flag.IntVar(&days, "days", 20, "Amount of days to retrieve time entries for")
	flag.IntVar(&customCapacity, "hours", 0, "Amount of hours in a single workweek (default: from harvest api)")
	flag.BoolVar(&onlyWorkedDays, "worked", false, "Only track days that have tracking entries")
	flag.StringVar(
		&group,
		"group",
		groupByDay,
		fmt.Sprintf(
			"Group results by %s|%s|%s|%s",
			groupByDay,
			groupByWeek,
			groupByMonth,
			groupByYear,
		),
	)
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

	switch group {
	case groupByDay:
	case groupByWeek:
	case groupByMonth:
	case groupByYear:
	default:
		return 1, fmt.Errorf("Invalid group '%s'", group)
	}

	from := time.Now()
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
		f, err := time.Parse(dateFormat, customDate)
		if err != nil {
			return 1, fmt.Errorf("Invalid date '%s' expected YYYY-mm-dd", customDate)
		}
		from = f
	}

	if err := t.SetUID(userID); err != nil {
		return 1, err
	}

	workWeek := float64(config.WorkWeek())
	capacity := Duration(t.User().Capacity())
	if customCapacity != 0 {
		capacity = Duration(customCapacity) * Duration(time.Hour)
	}
	daysCapacity := Duration(
		float64(capacity) * float64(days) / workWeek,
	)
	onlyWorkedDaysCopy := ""
	if onlyWorkedDays {
		onlyWorkedDaysCopy = " (estimate)"
	}

	c.l.Printf(
		"Running for %s %s\nID: %d\nWeek: %s\nOver %d days%s: %s\nFrom: %s\n\n",
		t.User().FirstName,
		t.User().LastName,
		t.User().ID,
		capacity,
		days,
		onlyWorkedDaysCopy,
		daysCapacity,
		from.Format("Mon Jan 02 2006"),
	)

	daysWorked, grouped, err := t.GetRecentDaysGrouped(days, from, !onlyWorkedDays, group)
	daysCapacity = Duration(
		float64(capacity) * float64(daysWorked) / workWeek,
	)

	var sum time.Duration
	for _, e := range grouped.SortSpent() {
		days := make(map[string]struct{}, 1)
		for _, d := range e.SpentDates {
			days[d.Format(dateFormat)] = struct{}{}
		}
		should := Duration(float64(capacity) * float64(len(days)) / workWeek)
		c.l.Printf(
			"%s - %5s / %s (%.2f%%)",
			e.FirstSpentDate.Format("Mon Jan 02 2006"),
			Duration(e.Hours),
			should,
			100*float64(e.Hours)/float64(should),
		)
		sum += e.Hours
	}

	diff := daysCapacity - Duration(sum)
	diffStr := fmt.Sprintf("%s remaining...", diff)
	if diff < 0 {
		diffStr = "target reached!"
	}

	c.l.Printf(
		"\nTotal: %s / %s (%.2f%%)\n%s",
		Duration(sum),
		daysCapacity,
		100*float64(sum)/float64(daysCapacity),
		diffStr,
	)

	return 0, nil
}
