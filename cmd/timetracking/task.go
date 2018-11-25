package main

import (
	"fmt"
	"sort"
	"strings"
)

type Task struct {
	ClientID    int    `json:"client_id"`
	ClientName  string `json:"client"`
	TaskID      int    `json:"task_id"`
	TaskName    string `json:"task"`
	ProjectID   int    `json:"project_id"`
	ProjectName string `json:"project"`
	ps          []string
}

func (t *Task) String() string {
	return fmt.Sprintf("%s %s %s", t.ClientName, t.ProjectName, t.TaskName)
}

type Tasks []*Task

type Score struct {
	Task  *Task
	Score int
}

type Scores []*Score

func (s Scores) Sort() Scores {
	sort.SliceStable(
		s,
		func(i, j int) bool {
			return s[i].Score > s[j].Score
		},
	)
	return s
}

func (s Scores) Tasks(amount int) Tasks {
	t := make(Tasks, 0, amount)
	for i := range s {
		if len(t) == amount {
			break
		}
		t = append(t, s[i].Task)
	}

	return t
}

func (t Tasks) FuzzyFind(find string, limit int, strict bool) Tasks {
	normalize := func(str string) string {
		return strings.ToLower(strings.Replace(str, " ", "", -1))
	}

	find = normalize(find)
	scores := make(Scores, 0, len(t))
	for _, t := range t {
		score := 0
		str := []rune(normalize(t.String()))
		start := 0
		runeCount := 0
	outer:
		for _, n := range find {
			runeCount++
			for f := start; f < len(str); f++ {
				if str[f] == n {
					start = f + 1
					score++
					continue outer
				}
			}
		}

		if strict && score < runeCount {
			continue
		}

		scores = append(scores, &Score{t, score})
	}

	return scores.Sort().Tasks(limit)
}
