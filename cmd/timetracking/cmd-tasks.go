package main

import (
	"flag"
)

func commandTasks(c *Command) (int, error) {
	var save bool
	flag.BoolVar(&save, "save", false, "Save in ~/.timetracking")
	flag.Parse()

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

	if err := t.SetUID(0); err != nil {
		return 1, err
	}

	res, err := t.GetUserProjectAssignments()
	d := make(Tasks, 0, len(res))
	if save {
		for _, a := range res {
			for _, t := range a.TaskAssignments {
				d = append(
					d,
					&Task{
						ClientID:    a.Client.ID,
						ProjectID:   a.Project.ID,
						TaskID:      t.Task.ID,
						ClientName:  a.Client.Name,
						ProjectName: a.Project.Name,
						TaskName:    t.Task.Name,
					},
				)
			}
		}

		config = &Config{}
		if err := confLoader.Read(config); err != nil {
			return 1, err
		}

		config.Tasks = d
		if err = confLoader.Create(config); err != nil {
			return 1, err
		}
		c.l.Println("Saved")
		return 0, nil
	}

	for _, a := range res {
		c.l.Printf("%s [%s]\n", a.Project.Name, a.Client.Name)
		for _, t := range a.TaskAssignments {
			c.l.Printf("    %s", t.Task.Name)
		}
	}

	return 0, nil
}
