package main

import (
	"fmt"
	"os"
	"strings"

	"regit/regit"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: re-git <command> [args]")
		return
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "init":
		regit.Init()
	case "add":
		for _, file := range args {
			regit.Add(file)
		}
	case "commit":
		if len(args) < 1 {
			fmt.Println("Commit message required")
			return
		}
		regit.Commit(strings.Join(args, " "))
	case "status":
		regit.Status()
	case "log":
		regit.Log()
	default:
		fmt.Println("Unknown command:", cmd)
	}
}
