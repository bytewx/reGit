package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"regit/regit"
)

func RunCLI() {
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
	case "remove":
		for _, file := range args {
			regit.Remove(file)
		}
	case "show":
		for _, file := range args {
			regit.Show(file)
		}
	case "ls-objects":
		regit.ListFiles()
	case "checkout":
		regit.Checkout()
	case "diff":
		regit.Diff()
	case "list-commits":
		regit.ListCommits()
	case "file-history":
		for _, file := range args {
			regit.FileHistory(file)
		}
	case "reset":
		regit.Reset()
	case "istracked":
		for _, file := range args {
			fmt.Println(file, regit.IsTracked(file))
		}
	case "get-file-version":
		if len(args) < 2 {
			fmt.Println("Usage: get-file-version <file> <commitIdx>")
			return
		}
		idx, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid commit index")
			return
		}
		regit.GetFileVersion(args[0], idx)
	case "commit-files":
		if len(args) < 1 {
			fmt.Println("Usage: commit-files <commitIdx>")
			return
		}
		idx, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid commit index")
			return
		}
		regit.CommitFiles(idx)
	case "remove-object":
		for _, oid := range args {
			regit.RemoveObject(oid)
		}
	case "commit-count":
		fmt.Println(regit.CommitCount())
	case "find-commit-by-message":
		for _, msg := range args {
			fmt.Println(regit.FindCommitByMessage(msg))
		}
	case "find-file-oids":
		for _, file := range args {
			fmt.Println(regit.FindFileOids(file))
		}
	case "restore-file-from-commit":
		if len(args) < 2 {
			fmt.Println("Usage: restore-file-from-commit <file> <commitIdx>")
			return
		}
		idx, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid commit index")
			return
		}
		regit.RestoreFileFromCommit(args[0], idx)
	case "purge-unreferenced-objects":
		regit.PurgeUnreferencedObjects()
	case "get-commit-message":
		for _, arg := range args {
			idx, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Invalid commit index")
				continue
			}
			fmt.Println(regit.GetCommitMessage(idx))
		}
	case "get-commit-date":
		for _, arg := range args {
			idx, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Invalid commit index")
				continue
			}
			fmt.Println(regit.GetCommitDate(idx))
		}
	case "get-commit-oid-for-file":
		if len(args) < 2 {
			fmt.Println("Usage: get-commit-oid-for-file <file> <commitIdx>")
			return
		}
		idx, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid commit index")
			return
		}
		fmt.Println(regit.GetCommitOidForFile(args[0], idx))
	case "list-all-tracked-files":
		files := regit.ListAllTrackedFiles()
		for _, f := range files {
			fmt.Println(f)
		}
	case "push":
		if len(args) < 1 {
			fmt.Println("Usage: push <remote_path>")
			return
		}
		regit.Push(args[0])
	default:
		fmt.Println("Unknown command:", cmd)
	}
}
