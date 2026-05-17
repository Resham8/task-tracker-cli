# task-cli

A fast, colourful task tracker for your terminal written in Go, powered by nothing but a JSON file.

```
$ task add "Write unit tests"
✓ Task #1 added: Write unit tests

$ task list
ID   Description                    Status      Created            Updated
──────────────────────────────────────────────────────────────────────────────────────────
1    Write unit tests               todo        17 May 2026        17 May 2026
──────────────────────────────────────────────────────────────────────────────────────────
 1 task(s)
```

---

## Features

- **Add, update, and delete** tasks from the command line
- **Three statuses** — `todo`, `in-progress`, `done` — with colour-coded output
- **Filter** the task list by status
- **Zero dependencies** — pure Go standard library
- **Persistent storage** via a local `tasks.json` file
- **Docker support** for a fully containerised workflow

---

## Requirements

- Go 1.21+ — [install here](https://go.dev/dl/)
- Or Docker, if you prefer containers

---

## Installation

**Clone and build**

```bash
git clone https://github.com/Resham8/task-tracker-cli.git
cd task-tracker-cli
go build -o task .
```

Move the binary somewhere on your `$PATH` for global access:

```bash
mv task /usr/local/bin/task
```


**Or use Docker**

```bash
docker build -t task-cli .
docker run --rm -v $(pwd):/app task-cli <command>
```

---

## Usage

```
task <command> [arguments]
```

| Command | Description |
|---|---|
| `add <description>` | Add a new task |
| `list [status]` | List all tasks, optionally filtered by status |
| `update <id> <description>` | Update a task's description |
| `status <id> <status>` | Change a task's status |
| `delete <id>` | Delete a task |
| `help`, `--help`, `-h` | Show help |

### Statuses

| Value | Meaning |
|---|---|
| `todo` | Not yet started (blue) |
| `in-progress` | Currently being worked on (yellow) |
| `done` | Completed (green) |

---

## Examples

```bash
# Add tasks
task add "Design the homepage"
task add "Write unit tests"
task add "Deploy to production"

# List everything
task list

# Filter by status
task list todo
task list in-progress
task list done

# Start working on task #1
task status 1 in-progress

# Finish it
task status 1 done

# Fix a typo in the description
task update 2 "Write integration tests"

# Remove a task
task delete 3
```

---

## Project Structure

```
task-tracker-cli/
├── main.go        # Entry point — argument parsing and command dispatch
├── commands.go    # Command implementations (add, list, update, status, delete)
├── task.go        # Task model, JSON persistence, and helper functions
├── task_test.go   # Tests
├── tasks.json     # Auto-generated local storage file
└── Dockerfile     # Container build
```

---

## Data Storage

Tasks are saved to `tasks.json` in the working directory. The file is created automatically on the first `add` command. Each task stores:

```json
[
  {
    "Id": 1,
    "Description": "Write unit tests",
    "Status": "in-progress",
    "CreatedAt": "17 May 2026 02:30 PM",
    "UpdatedAt": "17 May 2026 03:00 PM"
  }
]
```

---

## Running Tests

```bash
go test ./...
```

---

## License

MIT
