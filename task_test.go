package main

import (
	"path/filepath"
	"testing"
)

func setup(t *testing.T) {
	t.Helper()
	tasks = []Task{}
	dir := t.TempDir()

	storageFile = filepath.Join(dir, "temp_tasks.json")	
}

func TestCmdAdd(t *testing.T) {
	setup(t)
	cmdAdd([]string{"buy groceries"})

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	task := tasks[0]

	if task.Id != 1 {
		t.Errorf("Id = %d, want 1", task.Id)
	}

	if task.Description != "buy groceries" {
		t.Errorf("Description = %q, want %q", task.Description, "Buy groceries")
	}

	if task.Status != Todo {
		t.Errorf("Status = %q, want %q", task.Status, Todo)
	}

	if task.CreatedAt == "" || task.UpdatedAt == "" {
		t.Error("timestamps must not be empty")
	}
}

func TestCmdAddMultiWordArg(t *testing.T) {
	setup(t)
	cmdAdd([]string{"buy", "some", "pants"})

	if tasks[0].Description != "buy some pants" {
		t.Errorf("Description = %q, want %q", tasks[0].Description, "Buy more groceries")

	}

}

func TestCmdAddNoArgs(t *testing.T) {
	setup(t)
	cmdAdd(nil)

	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks after invalid add, got %d", len(tasks))
	}
}

func TestCmdAddIncrementsID(t *testing.T) {
	setup(t)
	cmdAdd([]string{"first"})
	cmdAdd([]string{"second"})
	cmdAdd([]string{"third"})

	for i, task := range tasks {
		want := uint(i + 1)
		if task.Id != want {
			t.Errorf("tasks[%d].Id = %d, want %d", i, task.Id, want)
		}
	}
}

// update tests

func TestCmdUpdate(t *testing.T) {
	setup(t)

	cmdAdd([]string{"original"})
	cmdUpdate([]string{"1", "update description"})

	if tasks[0].Description != "update description" {
		t.Errorf("Description = %q, want %q", tasks[0].Description, "update description")
	}

}

func TestCmdUpdateNotFound(t *testing.T) {
	setup(t)
	cmdAdd([]string{"task"})

	before := tasks[0].Description
	cmdUpdate([]string{"99", "new"})

	if tasks[0].Description != before {
		t.Error("description should be unchanged when ID not found")
	}
}

func TestCmdUpdateMissingArgs(t *testing.T) {
	setup(t)
	cmdAdd([]string{"task"})

	before := tasks[0].Description

	cmdUpdate([]string{"1"})

	if tasks[0].Description != before {
		t.Error("description should be unchanged with missing args")
	}
}

// status

func TestCmdStatus(t *testing.T) {
	setup(t)

	cmdAdd([]string{"task"})

	transitions := []Status{InProgress, Done, Todo}

	for _, s := range transitions {
		cmdStatus([]string{"1", string(s)})
		if tasks[0].Status != s {
			t.Errorf("after setting %q: Status = %q", s, tasks[0].Status)
		}
	}
}

func TestCmdStatusInvalid(t *testing.T) {
	setup(t)

	cmdAdd([]string{"task"})

	cmdStatus([]string{"1", "invalid"})

	if tasks[0].Status != Todo {
		t.Errorf("status should stay %q, got %q", Todo, tasks[0].Status)
	}
}

func TestCmdStatusNotFound(t *testing.T) {
	setup(t)
	cmdAdd([]string{"task"})
	cmdStatus([]string{"99", "done"})
	if tasks[0].Status != Todo {
		t.Errorf("status should be unchanged, got %q", tasks[0].Status)
	}
}

// Delete

func TestCmdDelete(t *testing.T) {
	setup(t)
	cmdAdd([]string{"a"})
	cmdAdd([]string{"b"})
	cmdAdd([]string{"c"})

	cmdDelete([]string{"2"})

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks after delete, got %d", len(tasks))

	}

	for _, t2 := range tasks {
		if t2.Id == 2 {
			t.Error("deleted task should not exist")
		}
	}
}

func TestCmdDeleteNotFound(t *testing.T) {
	setup(t)
	cmdAdd([]string{"task"})

	cmdDelete([]string{"99"})
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}

func TestCmdDeleteNoArgs(t *testing.T) {
	setup(t)
	cmdAdd([]string{"task"})

	cmdDelete(nil)
	if len(tasks) != 1 {

		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}

// cmdlist

func TestCmdListFilterByStatus(t *testing.T){
	setup(t)
	cmdAdd([]string{"todo task"})
	cmdAdd([]string{"done task"})
	cmdStatus([]string{"2", "done"})

	cmdList([]string{"done"})
	cmdList([]string{"todo"})
	cmdList([]string{"in-progress"})
	cmdList(nil)
}


func TestCmdListInvalidFilter(t *testing.T) {
	setup(t)
	cmdList([]string{"bogus"})
}