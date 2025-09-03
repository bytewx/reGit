package regit

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func Checkout() {
	commit := HeadCommit()
	if commit == "" {
		fmt.Println("No commits found")
		return
	}
	lines := strings.Split(commit, "\n")
	files := []string{}
	for _, line := range lines {
		if strings.Contains(line, " ") && !strings.HasPrefix(line, "commit") && !strings.HasPrefix(line, "Date:") {
			files = append(files, line)
		}
	}
	for _, entry := range files {
		parts := strings.Split(entry, " ")
		if len(parts) != 2 {
			continue
		}
		file, oid := parts[0], parts[1]
		objPath := filepath.Join(objectsDir, oid)
		data, err := ioutil.ReadFile(objPath)
		if err != nil {
			fmt.Println("Error restoring", file)
			continue
		}
		ioutil.WriteFile(file, data, 0644)
		fmt.Println("Restored", file)
	}
}

func Diff() {
	index, err := ioutil.ReadFile(indexFile)
	if err != nil {
		fmt.Println("Error reading index")
		return
	}
	lines := strings.Split(string(index), "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			continue
		}
		file, oid := parts[0], parts[1]
		objPath := filepath.Join(objectsDir, oid)
		staged, err := ioutil.ReadFile(objPath)
		if err != nil {
			continue
		}
		working, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("%s: file missing in working directory\n", file)
			continue
		}
		if string(staged) != string(working) {
			fmt.Printf("Diff for %s:\n", file)
			fmt.Println("--- staged")
			fmt.Println(string(staged))
			fmt.Println("--- working")
			fmt.Println(string(working))
		}
	}
}
