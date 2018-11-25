package main

import (
	"flag"
	"sort"
	"time"
)

func commandDaysOff(c *Command) (int, error) {
	var userID int
	var projectName string
	var hoursInt int
	var save bool
	flag.IntVar(&userID, "uid", 0, "The forecast user id of the user to fetch time-off entries for")
	flag.StringVar(&projectName, "project", "Time Off", "Name of the 'Time Off' project")
	flag.IntVar(&hoursInt, "hours", 7, "Amount of hours 'Time Off' should last for it to be an entire day off.")
	flag.BoolVar(&save, "save", false, "Save in ~/.timetracking")
	flag.Parse()

	hours := time.Hour * time.Duration(hoursInt)

	confLoader, config, err := getConfig(c.l)
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

	if err := t.SetForecastUID(userID); err != nil {
		return 1, err
	}

	r, err := t.GetAssignmentsByName(projectName)
	if err != nil {
		return 1, err
	}

	off := make([]time.Time, 0, len(r))
	for _, a := range r {
		if a.StartDate == nil || a.EndDate == nil {
			continue
		}

		s := a.StartDate.Time
		if !a.EndDate.After(s) {
			if a.Allocation.Duration >= hours {
				off = append(off, s)
			}
			continue
		}
		for {
			off = append(off, s)
			s = s.AddDate(0, 0, 1)
			if s.After(a.EndDate.Time) {
				break
			}
		}
	}

	sort.SliceStable(
		off,
		func(i, j int) bool { return off[i].Before(off[j]) },
	)

	if save {
		config = &Config{}
		if err := confLoader.Read(config); err != nil {
			return 1, err
		}
		unique := make(map[string]struct{}, len(off))
		for _, o := range off {
			unique[o.Format(dateFormat)] = struct{}{}
		}
		for _, o := range config.ExcludedDates {
			unique[o] = struct{}{}
		}

		uniqueSorted := make([]string, 0, len(unique))
		for i := range unique {
			uniqueSorted = append(uniqueSorted, i)
		}
		sort.SliceStable(
			uniqueSorted,
			func(i, j int) bool { return uniqueSorted[i] < uniqueSorted[j] },
		)

		config.ExcludedDates = uniqueSorted
		if err = confLoader.Create(config); err != nil {
			return 1, err
		}
		c.l.Println("Saved")
		return 0, nil
	}

	for _, d := range off {
		c.l.Println(d.Format(dateFormat))
	}

	return 0, nil
}
