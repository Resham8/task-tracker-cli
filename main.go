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

	if len(args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  add \"task\"")
		fmt.Println("  list")
		fmt.Println("  update <id> \"task\"")
		fmt.Println("  status <id> <todo|in-progress|done>")
		return
	}

	switch args[1] {
	case "list":
		listTasks()
	case "add":
		description := args[2]
		addTask(description)

	case "update":
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Invalid Id")
			return
		}

		description := args[3]
		updateTask(uint(id), description)

	case "status":
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Invalid Id")
			return
		}

		status := args[3]
		updateStatus(uint(id), Status(status))

	case "delete":
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Invalid Id")
			return
		}
		
		deleteTask(uint(id))

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
		Id:          getNextID(),
		Description: description,
		Status:      Todo,
		CreatedAt:   time.Now().Format("02 Jan 2006 03:04 PM"),
		UpdatedAt:   time.Now().Format("02 Jan 2006 03:04 PM"),
	}

	tasks = append(tasks, task)
	saveTasks()
	fmt.Println("Task added: ", tasks)
}

func updateTask(id uint, newDescription string) {
	for i, task := range tasks {
		if task.Id == id {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now().Format("02 Jan 2006 03:04 PM")
			saveTasks()
			fmt.Println("Task updated: ", tasks[i])
			return
		}
	}

	fmt.Println("Task not found")
}

func updateStatus(id uint, newStatus Status) {
	for i := range tasks {
		if tasks[i].Id == id {
			tasks[i].Status = newStatus
			tasks[i].UpdatedAt = time.Now().Format("02 Jan 2006 03:04 PM")

			saveTasks()
			fmt.Println("Task marked as ", newStatus)
			return
		}
	}
}

func listTasks() {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	for _, t := range tasks {
		fmt.Printf("[%d] %s\n", t.Id, t.Description)
		fmt.Printf("   Status: %s\n", t.Status)
		fmt.Printf("   Created: %s\n", t.CreatedAt)
		fmt.Printf("   Updated: %s\n\n", t.UpdatedAt)
	}
}

func deleteTask(id uint) {
	for i := range tasks {
		if tasks[i].Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			saveTasks()
			fmt.Println("Task deleted")
			return
		}
	}

	fmt.Println("Task not found")
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
