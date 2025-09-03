package regit

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func Log() {
	log, _ := ioutil.ReadFile(logFile)
	entries := strings.Split(string(log), "---\n")
	for _, entry := range entries {
		if strings.TrimSpace(entry) != "" {
			fmt.Println(entry)
		}
	}
}

func ListCommits() {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Println("Error reading log")
		return
	}
	entries := strings.Split(string(log), "---\n")
	for _, entry := range entries {
		lines := strings.Split(entry, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "commit ") || strings.HasPrefix(line, "Date:") {
				fmt.Println(line)
			}
			if strings.HasPrefix(line, "commit ") {
				fmt.Println("-----")
			}
		}
	}
}

func FileHistory(file string) {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Println("Error reading log")
		return
	}
	entries := strings.Split(string(log), "---\n")
	for _, entry := range entries {
		if strings.Contains(entry, file+" ") {
			lines := strings.Split(entry, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "commit ") || strings.HasPrefix(line, "Date:") || strings.HasPrefix(line, file+" ") {
					fmt.Println(line)
				}
			}
			fmt.Println("-----")
		}
	}
}

func CommitCount() int {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		return 0
	}
	entries := strings.Split(string(log), "---\n")
	return len(entries) - 1
}

func FindCommitByMessage(substring string) []int {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Println("Error reading log")
		return nil
	}
	entries := strings.Split(string(log), "---\n")
	var indices []int
	for i, entry := range entries[:len(entries)-1] {
		if strings.Contains(entry, substring) {
			indices = append(indices, i)
		}
	}
	return indices
}

func FindFileOids(file string) []string {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Println("Error reading log")
		return nil
	}
	entries := strings.Split(string(log), "---\n")
	var oids []string
	for _, entry := range entries[:len(entries)-1] {
		lines := strings.Split(entry, "\n")
		for _, line := range lines {
			parts := strings.Split(line, " ")
			if len(parts) == 2 && parts[0] == file {
				oids = append(oids, parts[1])
			}
		}
	}
	return oids
}

func ListAllTrackedFiles() []string {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		return nil
	}
	entries := strings.Split(string(log), "---\n")
	filesSet := make(map[string]struct{})
	for _, entry := range entries[:len(entries)-1] {
		lines := strings.Split(entry, "\n")
		for _, line := range lines {
			parts := strings.Split(line, " ")
			if len(parts) == 2 {
				filesSet[parts[0]] = struct{}{}
			}
		}
	}
	files := make([]string, 0, len(filesSet))
	for f := range filesSet {
		files = append(files, f)
	}
	return files
}

func GetCommitMessage(commitIdx int) string {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		return ""
	}
	entries := strings.Split(string(log), "---\n")
	if commitIdx < 0 || commitIdx >= len(entries)-1 {
		return ""
	}
	entry := entries[commitIdx]
	lines := strings.Split(entry, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "commit ") {
			for j := i + 1; j < len(lines); j++ {
				if strings.HasPrefix(lines[j], "Date:") && j+1 < len(lines) {
					return lines[j+1]
				}
			}
		}
	}
	return ""
}

func GetCommitDate(commitIdx int) string {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		return ""
	}
	entries := strings.Split(string(log), "---\n")
	if commitIdx < 0 || commitIdx >= len(entries)-1 {
		return ""
	}
	entry := entries[commitIdx]
	lines := strings.Split(entry, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Date:") {
			return strings.TrimPrefix(line, "Date: ")
		}
	}
	return ""
}

func GetCommitOidForFile(file string, commitIdx int) string {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		return ""
	}
	entries := strings.Split(string(log), "---\n")
	if commitIdx < 0 || commitIdx >= len(entries)-1 {
		return ""
	}
	entry := entries[commitIdx]
	lines := strings.Split(entry, "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 2 && parts[0] == file {
			return parts[1]
		}
	}
	return ""
}

func Blame(file string) {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Println("Error reading log")
		return
	}
	entries := strings.Split(string(log), "---\n")
	lineCommit := make(map[int]int)
	var fileLines []string
	for idx, entry := range entries[:len(entries)-1] {
		lines := strings.Split(entry, "\n")
		var oid string
		for _, line := range lines {
			parts := strings.Split(line, " ")
			if len(parts) == 2 && parts[0] == file {
				oid = parts[1]
				break
			}
		}
		if oid != "" {
			objPath := filepath.Join(objectsDir, oid)
			data, err := ioutil.ReadFile(objPath)
			if err == nil {
				fileLines = strings.Split(string(data), "\n")
				for i := range fileLines {
					lineCommit[i] = idx
				}
			}
		}
	}
	for i, line := range fileLines {
		fmt.Printf("%d %s | %s\n", lineCommit[i], GetCommitMessage(lineCommit[i]), line)
	}
}

func ShowCommitFiles(commitIdx int) {
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
	fmt.Printf("Files in commit %d:\n", commitIdx)
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			fmt.Println(parts[0])
		}
	}
}

func ShowCommitDiff(file string, idxA, idxB int) {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Println("Error reading log")
		return
	}
	entries := strings.Split(string(log), "---\n")
	if idxA < 0 || idxA >= len(entries)-1 || idxB < 0 || idxB >= len(entries)-1 {
		fmt.Println("Invalid commit index")
		return
	}
	getOid := func(entry string) string {
		lines := strings.Split(entry, "\n")
		for _, line := range lines {
			parts := strings.Split(line, " ")
			if len(parts) == 2 && parts[0] == file {
				return parts[1]
			}
		}
		return ""
	}
	oidA := getOid(entries[idxA])
	oidB := getOid(entries[idxB])
	if oidA == "" || oidB == "" {
		fmt.Println("File not found in one of the commits")
		return
	}
	dataA, errA := ioutil.ReadFile(filepath.Join(objectsDir, oidA))
	dataB, errB := ioutil.ReadFile(filepath.Join(objectsDir, oidB))
	if errA != nil || errB != nil {
		fmt.Println("Error reading file objects")
		return
	}
	fmt.Printf("Diff for %s between commit %d and %d:\n", file, idxA, idxB)
	fmt.Println("--- commit", idxA)
	fmt.Println(string(dataA))
	fmt.Println("--- commit", idxB)
	fmt.Println(string(dataB))
}
