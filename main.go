package main

import (
	"fmt"
	"os"
)

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
