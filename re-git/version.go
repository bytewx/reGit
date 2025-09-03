package regit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func GetFileVersion(file string, commitIdx int) {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Println("Error reading log")
		return
	}
	entries := strings.Split(string(log), "---\n")
	if commitIdx < 0 || commitIdx >= len(entries)-1 {
		fmt.Println("Invalid commit index")
		return
	}
	entry := entries[commitIdx]
	lines := strings.Split(entry, "\n")
	var oid string
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 2 && parts[0] == file {
			oid = parts[1]
			break
		}
	}
	if oid == "" {
		fmt.Println("File not found in commit")
		return
	}
	objPath := filepath.Join(objectsDir, oid)
	data, err := ioutil.ReadFile(objPath)
	if err != nil {
		fmt.Println("Object not found")
		return
	}
	fmt.Printf("Version of %s from commit %d:\n%s\n", file, commitIdx, string(data))
}

func RestoreFileFromCommit(file string, commitIdx int) {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Println("Error reading log")
		return
	}
	entries := strings.Split(string(log), "---\n")
	if commitIdx < 0 || commitIdx >= len(entries)-1 {
		fmt.Println("Invalid commit index")
		return
	}
	entry := entries[commitIdx]
	lines := strings.Split(entry, "\n")
	var oid string
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 2 && parts[0] == file {
			oid = parts[1]
			break
		}
	}
	if oid == "" {
		fmt.Println("File not found in commit")
		return
	}
	objPath := filepath.Join(objectsDir, oid)
	data, err := ioutil.ReadFile(objPath)
	if err != nil {
		fmt.Println("Object not found")
		return
	}
	err = ioutil.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Println("Error restoring file")
		return
	}
	fmt.Printf("Restored %s from commit %d\n", file, commitIdx)
}

func Revert(commitIdx int) {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Println("Error reading log")
		return
	}
	entries := strings.Split(string(log), "---\n")
	if commitIdx < 0 || commitIdx >= len(entries)-1 {
		fmt.Println("Invalid commit index")
		return
	}
	entry := entries[commitIdx]
	lines := strings.Split(entry, "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			os.Remove(parts[0])
		}
	}
	fmt.Println("Reverted commit", commitIdx)
}

func CherryPick(commitIdx int) {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Println("Error reading log")
		return
	}
	entries := strings.Split(string(log), "---\n")
	if commitIdx < 0 || commitIdx >= len(entries)-1 {
		fmt.Println("Invalid commit index")
		return
	}
	entry := entries[commitIdx]
	lines := strings.Split(entry, "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			objPath := filepath.Join(objectsDir, parts[1])
			data, err := ioutil.ReadFile(objPath)
			if err == nil {
				ioutil.WriteFile(parts[0], data, 0644)
			}
		}
	}
	fmt.Println("Cherry-picked commit", commitIdx)
}
