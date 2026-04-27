package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
	red    = "\033[31m"
	gray   = "\033[90m"
)

func colorize(s, color string) string { return color + s + reset }

var tasks = []Task{}
var validStatuses = []Status{Todo, InProgress, Done}

const storageFile = "tasks.json"

func (s Status) isValid() bool {
	for _, v := range validStatuses {
		if s == v {
			return true
		}
	}
	return false
}

func (s Status) colored() string {
	switch s {
	case Todo:
		return colorize(string(s), blue)
	case InProgress:
		return colorize(string(s), yellow)
	case Done:
		return colorize(string(s), green)
	default:
		return string(s)
	}
}

func printHelp() {
	fmt.Println(colorize(bold+"Task CLI"+reset, cyan) + "  — manage tasks from your terminal\n")
	cmds := [][2]string{
		{"add <description>", "Add a new task"},
		{"list [status]", "List tasks (optionally filter by status)"},
		{"update <id> <description>", "Update a task's description"},
		{"status <id> <status>", "Change a task's status"},
		{"delete <id>", "Delete a task"},
	}
	fmt.Println(colorize("Usage:", bold))
	for _, c := range cmds {
		fmt.Printf("  %-35s %s\n", colorize(c[0], cyan), c[1])
	}
	fmt.Println("\n" + colorize("Statuses:", bold))
	fmt.Printf("  %s  %s  %s\n",
		Todo.colored(), InProgress.colored(), Done.colored())
	fmt.Println("\n" + colorize("Examples:", bold))
	fmt.Println("  task add \"Write unit tests\"")
	fmt.Println("  task list in-progress")
	fmt.Println("  task status 1 done")
}

func main() {
	if err := loadTasks(); err != nil {
		printError("failed to load tasks: " + err.Error())
		os.Exit(1)
	}

	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		return
	}

	cmd, rest := args[0], args[1:]

	switch cmd {
	case "add":
		cmdAdd(rest)
	case "list":
		cmdList(rest)
	case "update":
		cmdUpdate(rest)
	case "status":
		cmdStatus(rest)
	case "delete":
		cmdDelete(rest)
	case "help", "--help", "-h":
		printHelp()
	default:
		printError(fmt.Sprintf("unknown command %q", cmd))
		fmt.Println()
		printHelp()

	}

}

func printError(msg string)   { fmt.Fprintln(os.Stderr, colorize("✗ "+msg, red)) }
func printSuccess(msg string) { fmt.Println(colorize("✓ ", green) + msg) }

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

func cmdAdd(args []string) {
	if len(args) == 0 {
		printError("missing description — usage: add \"<description>\"")
		return
	}
	description := strings.Join(args, " ")
	task := Task{
		Id:          getNextID(),
		Description: description,
		Status:      Todo,
		CreatedAt:   time.Now().Format("02 Jan 2006 03:04 PM"),
		UpdatedAt:   time.Now().Format("02 Jan 2006 03:04 PM"),
	}

	tasks = append(tasks, task)
	if err := saveTasks(); err != nil {
		printError(err.Error())
		return
	}
	printSuccess(fmt.Sprintf("Task #%d added: %s", task.Id, colorize(task.Description, cyan)))
}

func cmdUpdate(args []string) {
	if len(args) < 2 {
		printError("usage: update <id> \"<new description>\"")
		return
	}
	id, err := parseID(args[0])
	if err != nil {
		printError(err.Error())
		return
	}
	newDesc := strings.Join(args[1:], " ")
	for i, task := range tasks {
		if task.Id == id {
			tasks[i].Description = newDesc
			tasks[i].UpdatedAt = time.Now().Format("02 Jan 2006 03:04 PM")
			if err := saveTasks(); err != nil {
				printError(err.Error())
				return
			}
			printSuccess(fmt.Sprintf("Task #%d updated: %s", id, colorize(newDesc, cyan)))
			return
		}
	}

	printError(fmt.Sprintf("no task with id %d", id))
}

func cmdStatus(args []string) {
	if len(args) < 2 {
		printError("usage: status <id> <todo|in-progress|done>")
		return
	}
	id, err := parseID(args[0])
	if err != nil {
		printError(err.Error())
		return
	}
	newStatus := Status(args[1])

	if !newStatus.isValid() {
		printError(fmt.Sprintf("invalid status %q — choose: todo, in-progress, done", newStatus))
		return
	}

	for i := range tasks {
		if tasks[i].Id == id {
			tasks[i].Status = newStatus
			tasks[i].UpdatedAt = time.Now().Format("02 Jan 2006 03:04 PM")

			if err := saveTasks(); err != nil {
				printError(err.Error())
				return
			}
			printSuccess(fmt.Sprintf("Task #%d marked as %s", id, newStatus.colored()))
			return
		}
	}

	printError(fmt.Sprintf("no task with id %d", id))
}

func cmdList(args []string) {
	var filter Status
	if len(args) > 0 {
		filter = Status(args[0])
		if !filter.isValid() {
			printError(fmt.Sprintf("unknown status %q — choose: todo, in-progress, done", filter))
			return
		}
	}

	visible := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		if filter == "" || t.Status == filter {
			visible = append(visible, t)
		}
	}

	if len(visible) == 0 {
		if filter != "" {
			fmt.Printf("%s  No tasks with status %s\n", colorize("—", gray), filter.colored())
		} else {
			fmt.Println(colorize("—", gray) + "  No tasks yet. Try: add \"your first task\"")
		}
		return
	}

	header := fmt.Sprintf("%-4s  %-30s  %-11s  %-18s  %-18s",
		"ID", "Description", "Status", "Created", "Updated")
	divider := strings.Repeat("─", 90)

	fmt.Println(colorize(bold+header, gray))
	fmt.Println(colorize(divider, gray))

	for _, t := range visible {
		fmt.Printf("%-4d  %-30s  %-20s  %-18s  %-18s\n",
			t.Id,
			truncate(t.Description, 30),
			t.Status.colored(),
			t.CreatedAt,
			t.UpdatedAt,
		)
	}

	fmt.Println(colorize(divider, gray))
	fmt.Printf(dim+"  %d task(s)"+reset+"\n", len(visible))
}

func cmdDelete(args []string) {
	if len(args) == 0 {
		printError("usage: delete <id>")
		return
	}
	id, err := parseID(args[0])
	if err != nil {
		printError(err.Error())
		return
	}
	for i := range tasks {
		if tasks[i].Id == id {
			desc := tasks[i].Description
			tasks = append(tasks[:i], tasks[i+1:]...)
			if err := saveTasks(); err != nil {
				printError(err.Error())
				return
			}
			printSuccess(fmt.Sprintf("Task #%d deleted: %s", id, colorize(desc, gray)))
			return
		}
	}
	printError(fmt.Sprintf("no task with id %d", id))
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
