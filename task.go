package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Status string

const (
	Todo       Status = "todo"
	InProgress Status = "in-progress"
	Done       Status = "done"
)

type Task struct {
	Id          uint
	Description string
	Status      Status
	CreatedAt   string
	UpdatedAt   string
}

var tasks = []Task{}
var validStatuses = []Status{Todo, InProgress, Done}

var storageFile = "tasks.json"

func (s Status) isValid() bool {
	for _, v := range validStatuses {
		if s == v {
			return true
		}
	}
	return false
}

func parseID(raw string) (uint, error) {
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("invalid id %q — must be a positive number", raw)
	}
	return uint(n), nil
}

func loadTasks() error {
	data, err := os.ReadFile(storageFile)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", storageFile, err)
	}
	return json.Unmarshal(data, &tasks)
}

func saveTasks() error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding tasks: %w", err)
	}
	if err := os.WriteFile(storageFile, data, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", storageFile, err)
	}

	return nil
}

func getNextID() uint {
	var max uint
	for _, t := range tasks {
		if t.Id > max {
			max = t.Id
		}
	}
	return max + 1
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
