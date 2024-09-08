package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const MAX_DAYS = 180
const MAX_WEEKS = 25
const OUT_OF_RANGE = 9999

// Read commits from saved repos and print graph
func graph(email string) {
	fmt.Printf("👤 GitGraph User: %s\n", email)
	commits, err := processRepos(email)
	if err != nil {
		panic(err)
	}
	printGraph(commits)
}

// Print graph from commits map
func printGraph(commits map[int]int) {
	keys := getSortedKeys(commits)
	graph := makeGraphMap(keys, commits)
	printCells(graph)
}

// Get the keys of a map in sorted order
func getSortedKeys(m map[int]int) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	return keys
}

type column []int

// Build a graph map with weeks as keys and commit count array as values
func makeGraphMap(keys []int, commits map[int]int) map[int]column {
	weeks := make(map[int]column)
	col := column{}

	for _, k := range keys {
		week := int(k / 7)
		currDay := k % 7

		if currDay == 0 {
			col = column{}
		}

		col = append(col, commits[k])

		if currDay == 6 {
			weeks[week] = col
		}
	}

	return weeks
}

// Print the graph
func printCells(weeks map[int]column) {
	printMonths()
	for j := 6; j >= 0; j-- {
		for i := MAX_WEEKS + 1; i >= 0; i-- {
			if i == MAX_WEEKS+1 {
				printDayCol(j)
			}

			if week, ok := weeks[i]; ok {
				// Today
				if i == 0 && j == getMissingDays()-1 {
					printCell(week[j], true)
					continue
				} else {
					if j < len(week) {
						printCell(week[j], false)
						continue
					}
				}
			}

			printCell(0, false)
		}
		fmt.Println()
	}
}

// Print month names
func printMonths() {
	week := resetTime(time.Now()).Add(-(MAX_DAYS * time.Hour * 24))
	month := week.Month()
	fmt.Printf("             ")

	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}
		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Println()
}

// Print cell with bgcolored by value
func printCell(val int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	case val >= 10:
		escape = "\033[1;30;42m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

// Print day names
func printDayCol(day int) {
	out := "     "
	switch day {
	case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri "
	}

	fmt.Print(out)
}

// Read saved repos and get commits map
func processRepos(email string) (map[int]int, error) {
	dotFile := getDotFilePath()
	repos := readRepoOnly(dotFile)

	if len(repos) == 0 {
		fmt.Println("🔍 Please scan git repos with '-scan folder'")
		panic("No repos found")
	}

	totalDays := MAX_DAYS

	commits := make(map[int]int, totalDays)
	for i := 0; i < totalDays; i++ {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = getCommits(email, path, commits)
	}

	return commits, nil
}

// Get commits map from a repo
func getCommits(email string, path string, commits map[int]int) map[int]int {
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}

	itr, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}

	offset := getMissingDays()
	err = itr.ForEach(func(c *object.Commit) error {
		daysAgo := countDaysSince(c.Author.When) + offset

		if c.Author.Email != email {
			return nil
		}

		if daysAgo != OUT_OF_RANGE {
			commits[daysAgo]++
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return commits
}

// Count how many days ago from current time
func countDaysSince(t time.Time) int {
	days := 0
	now := resetTime(time.Now())

	for t.Before(now) {
		t = t.Add(time.Hour * 24)
		days++
		if days > MAX_DAYS {
			return OUT_OF_RANGE
		}
	}

	return days
}

// Set time of a date to 0
func resetTime(t time.Time) time.Time {
	year, month, day := t.Date()
	startTime := time.Date(year, month, day, 0, 0, 0, 0, t.Location())

	return startTime
}

// Get missing days to fill in the graph
func getMissingDays() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}

	return offset
}
