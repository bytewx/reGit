package regit

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	repoDir    = ".git"
	objectsDir = ".git/objects"
	indexFile  = ".git/index"
	logFile    = ".git/log"
	headFile   = ".git/HEAD"
	refsDir    = ".git/refs"
	headsDir   = ".git/refs/heads"
)

var stashFile = ".git/stash"

func Init() {
	os.Mkdir(repoDir, 0755)
	os.Mkdir(objectsDir, 0755)
	os.Mkdir(refsDir, 0755)
	os.Mkdir(headsDir, 0755)
	ioutil.WriteFile(indexFile, []byte{}, 0644)
	ioutil.WriteFile(logFile, []byte{}, 0644)
	ioutil.WriteFile(headFile, []byte("ref: refs/heads/master\n"), 0644)
	ioutil.WriteFile(filepath.Join(headsDir, "master"), []byte{}, 0644)
	fmt.Println("Initialized empty re-git repository in", repoDir)
}

func Add(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading file:", file)
		return
	}
	hash := sha1.Sum(data)
	oid := hex.EncodeToString(hash[:])
	objPath := filepath.Join(objectsDir, oid)
	ioutil.WriteFile(objPath, data, 0644)

	index, _ := ioutil.ReadFile(indexFile)
	lines := strings.Split(string(index), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, file+" ") {
			fmt.Println(file, "already staged")
			return
		}
	}
	f, _ := os.OpenFile(indexFile, os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(fmt.Sprintf("%s %s\n", file, oid))
	fmt.Println("Added", file)
}

func Commit(message string) {
	index, err := ioutil.ReadFile(indexFile)
	if err != nil || len(index) == 0 {
		fmt.Println("Nothing to commit")
		return
	}
	timestamp := time.Now().Format(time.RFC3339)
	entry := fmt.Sprintf("commit %s\nDate: %s\n\n%s\n%s\n---\n", timestamp, timestamp, message, string(index))
	f, _ := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(entry)
	ioutil.WriteFile(indexFile, []byte{}, 0644)
	fmt.Println("Committed:", message)
}

func Status() {
	index, _ := ioutil.ReadFile(indexFile)
	if len(index) == 0 {
		fmt.Println("No files staged")
		return
	}
	fmt.Println("Staged files:")
	lines := strings.Split(string(index), "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Println(" ", strings.Split(line, " ")[0])
		}
	}
}

func Remove(file string) {
	index, err := ioutil.ReadFile(indexFile)
	if err != nil {
		fmt.Println("Error reading index")
		return
	}
	lines := strings.Split(string(index), "\n")
	newLines := []string{}
	removed := false
	for _, line := range lines {
		if strings.HasPrefix(line, file+" ") {
			removed = true
			continue
		}
		if line != "" {
			newLines = append(newLines, line)
		}
	}
	if removed {
		ioutil.WriteFile(indexFile, []byte(strings.Join(newLines, "\n")+"\n"), 0644)
		fmt.Println("Removed", file, "from staging")
	} else {
		fmt.Println(file, "not staged")
	}
}

func Reset() {
	err := ioutil.WriteFile(indexFile, []byte{}, 0644)
	if err != nil {
		fmt.Println("Error resetting index")
		return
	}
	fmt.Println("Staging area cleared")
}

func CreateBranch(name string) {
	branchPath := filepath.Join(headsDir, name)
	if _, err := os.Stat(branchPath); err == nil {
		fmt.Println("Branch already exists:", name)
		return
	}
	err := ioutil.WriteFile(branchPath, []byte{}, 0644)
	if err != nil {
		fmt.Println("Error creating branch:", name)
		return
	}
	fmt.Println("Created branch:", name)
}

func ListBranches() {
	files, err := ioutil.ReadDir(headsDir)
	if err != nil {
		fmt.Println("Error reading branches")
		return
	}
	fmt.Println("Branches:")
	for _, f := range files {
		fmt.Println(" ", f.Name())
	}
}

func StashSave() {
	index, err := ioutil.ReadFile(indexFile)
	if err != nil || len(index) == 0 {
		fmt.Println("Nothing to stash")
		return
	}
	err = ioutil.WriteFile(stashFile, index, 0644)
	if err != nil {
		fmt.Println("Error saving stash")
		return
	}
	ioutil.WriteFile(indexFile, []byte{}, 0644)
	fmt.Println("Stashed current staged files")
}

func StashApply() {
	stash, err := ioutil.ReadFile(stashFile)
	if err != nil || len(stash) == 0 {
		fmt.Println("No stash found")
		return
	}
	err = ioutil.WriteFile(indexFile, stash, 0644)
	if err != nil {
		fmt.Println("Error applying stash")
		return
	}
	fmt.Println("Applied stash to staging area")
}

func StashDrop() {
	err := ioutil.WriteFile(stashFile, []byte{}, 0644)
	if err != nil {
		fmt.Println("Error dropping stash")
		return
	}
	fmt.Println("Dropped stash")
}

func Rename(oldName, newName string) {
	index, err := ioutil.ReadFile(indexFile)
	if err != nil {
		fmt.Println("Error reading index")
		return
	}
	lines := strings.Split(string(index), "\n")
	updated := false
	for i, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 2 && parts[0] == oldName {
			lines[i] = newName + " " + parts[1]
			updated = true
			break
		}
	}
	if updated {
		err := os.Rename(oldName, newName)
		if err != nil {
			fmt.Println("Error renaming file in working directory")
			return
		}
		ioutil.WriteFile(indexFile, []byte(strings.Join(lines, "\n")), 0644)
		fmt.Printf("Renamed %s to %s\n", oldName, newName)
	} else {
		fmt.Println("File not staged:", oldName)
	}
}

func Move(file, newDir string) {
	newPath := filepath.Join(newDir, filepath.Base(file))
	err := os.Rename(file, newPath)
	if err != nil {
		fmt.Println("Error moving file")
		return
	}
	index, err := ioutil.ReadFile(indexFile)
	if err != nil {
		fmt.Println("Error reading index")
		return
	}
	lines := strings.Split(string(index), "\n")
	for i, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 2 && parts[0] == file {
			lines[i] = newPath + " " + parts[1]
			break
		}
	}
	ioutil.WriteFile(indexFile, []byte(strings.Join(lines, "\n")), 0644)
	fmt.Printf("Moved %s to %s\n", file, newPath)
}
