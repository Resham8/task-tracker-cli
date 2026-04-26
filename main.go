package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
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

func main() {
	args := os.Args

	loadTasks()
	if len(args) < 3 {
		fmt.Println("Usage: go run main.go add \"task description\"")
		return
	}

	switch args[1] {
	case "add":
		description := args[2]
		addTask(description)

	case "update":
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Invalid Id")
			return
		}

		description  := args[3]
		updateTask(uint(id),description)
	}

}
func loadTasks() error {
	data, err := os.ReadFile("tasks.json")
	if err != nil {
		return nil
	}

	return json.Unmarshal(data, &tasks)
}

func saveTasks() error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile("tasks.json", data, 0644)
}

func addTask(description string) {
	task := Task{
		Id:          uint(len(tasks) + 1),
		Description: description,
		Status:      Todo,
		CreatedAt:   time.Now().Format("02 Jan 2006 03:04 PM"),
		UpdatedAt:   time.Now().Format("02 Jan 2006 03:04 PM"),
	}

	tasks = append(tasks, task)
	saveTasks()
	fmt.Println("Task added: ", tasks)
}

func updateTask(id uint, description string) {

}
