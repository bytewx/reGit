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
	case "help":
		fmt.Println(`Available commands:
			init
			add <file>
			commit "<message>"
			status
			log
			remove <file>
			show <file>
			ls-objects
			checkout
			diff
			list-commits
			file-history <file>
			reset
			istracked <file>
			get-file-version <file> <commitIdx>
			commit-files <commitIdx>
			remove-object <oid>
			commit-count
			find-commit-by-message "<msg>"
			find-file-oids <file>
			restore-file-from-commit <file> <commitIdx>
			purge-unreferenced-objects
			get-commit-message <commitIdx>
			get-commit-date <commitIdx>
			get-commit-oid-for-file <file> <commitIdx>
			list-all-tracked-files
			push <remote_path>
			pull <remote_path>
			clone <remote_path> <target_path>
			fetch <remote_path>
			merge <remote_path>
			merge-to-remote <remote_path>
			help`)
		return
	case "stash-save":
		regit.StashSave()
	case "stash-apply":
		regit.StashApply()
	case "stash-drop":
		regit.StashDrop()
	case "blame":
		for _, file := range args {
			regit.Blame(file)
		}
	case "revert":
		for _, arg := range args {
			idx, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Invalid commit index")
				continue
			}
			regit.Revert(idx)
		}
	case "cherry-pick":
		for _, arg := range args {
			idx, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Invalid commit index")
				continue
			}
			regit.CherryPick(idx)
		}
	case "rename":
		if len(args) < 2 {
			fmt.Println("Usage: rename <oldName> <newName>")
			return
		}
		regit.Rename(args[0], args[1])
	case "move":
		if len(args) < 2 {
			fmt.Println("Usage: move <file> <newDir>")
			return
		}
		regit.Move(args[0], args[1])
	case "show-commit-files":
		for _, arg := range args {
			idx, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Invalid commit index")
				continue
			}
			regit.ShowCommitFiles(idx)
		}
	case "show-commit-diff":
		if len(args) < 3 {
			fmt.Println("Usage: show-commit-diff <file> <commitIdxA> <commitIdxB>")
			return
		}
		idxA, errA := strconv.Atoi(args[1])
		idxB, errB := strconv.Atoi(args[2])
		if errA != nil || errB != nil {
			fmt.Println("Invalid commit index")
			return
		}
		regit.ShowCommitDiff(args[0], idxA, idxB)
	default:
		fmt.Println("Unknown command:", cmd)
	}
}
