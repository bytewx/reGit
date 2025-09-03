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
	repoDir      = ".git"
	objectsDir   = ".git/objects"
	indexFile    = ".git/index"
	logFile      = ".git/log"
	headFile     = ".git/HEAD"
	refsDir      = ".git/refs"
	headsDir     = ".git/refs/heads"
	tagsDir      = ".git/refs/tags"
	configFile   = ".git/config"
)

func Init() {
	os.Mkdir(repoDir, 0755)
	os.Mkdir(objectsDir, 0755)
	ioutil.WriteFile(indexFile, []byte{}, 0644)
	ioutil.WriteFile(logFile, []byte{}, 0644)
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

func Log() {
	log, _ := ioutil.ReadFile(logFile)
	entries := strings.Split(string(log), "---\n")
	for _, entry := range entries {
		if strings.TrimSpace(entry) != "" {
			fmt.Println(entry)
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

func Show(file string) {
	index, err := ioutil.ReadFile(indexFile)
	if err != nil {
		fmt.Println("Error reading index")
		return
	}
	lines := strings.Split(string(index), "\n")
	var oid string
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 2 && parts[0] == file {
			oid = parts[1]
			break
		}
	}
	if oid == "" {
		fmt.Println(file, "not staged")
		return
	}
	objPath := filepath.Join(objectsDir, oid)
	data, err := ioutil.ReadFile(objPath)
	if err != nil {
		fmt.Println("Object not found")
		return
	}
	fmt.Printf("Contents of %s:\n%s\n", file, string(data))
}

func ListFiles() {
	files, err := ioutil.ReadDir(objectsDir)
	if err != nil {
		fmt.Println("Error reading objects")
		return
	}
	fmt.Println("Tracked objects:")
	for _, f := range files {
		fmt.Println(" ", f.Name())
	}
}

func HeadCommit() string {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		return ""
	}
	entries := strings.Split(string(log), "---\n")
	if len(entries) == 0 {
		return ""
	}
	last := strings.TrimSpace(entries[len(entries)-2])
	return last
}

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

func Reset() {
	err := ioutil.WriteFile(indexFile, []byte{}, 0644)
	if err != nil {
		fmt.Println("Error resetting index")
		return
	}
	fmt.Println("Staging area cleared")
}

func IsTracked(file string) bool {
	files, err := ioutil.ReadDir(objectsDir)
	if err != nil {
		return false
	}
	for _, f := range files {
		objPath := filepath.Join(objectsDir, f.Name())
		data, err := ioutil.ReadFile(objPath)
		if err != nil {
			continue
		}
		working, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		if string(data) == string(working) {
			return true
		}
	}
	return false
}

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

func CommitFiles(commitIdx int) {
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

func RemoveObject(oid string) {
	objPath := filepath.Join(objectsDir, oid)
	err := os.Remove(objPath)
	if err != nil {
		fmt.Println("Error removing object:", oid)
		return
	}
	fmt.Println("Removed object:", oid)
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

func PurgeUnreferencedObjects() {
	log, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Println("Error reading log")
		return
	}
	entries := strings.Split(string(log), "---\n")
	referenced := make(map[string]bool)
	for _, entry := range entries[:len(entries)-1] {
		lines := strings.Split(entry, "\n")
		for _, line := range lines {
			parts := strings.Split(line, " ")
			if len(parts) == 2 {
				referenced[parts[1]] = true
			}
		}
	}
	files, err := ioutil.ReadDir(objectsDir)
	if err != nil {
		fmt.Println("Error reading objects")
		return
	}
	for _, f := range files {
		if !referenced[f.Name()] {
			objPath := filepath.Join(objectsDir, f.Name())
			os.Remove(objPath)
			fmt.Println("Purged unreferenced object:", f.Name())
		}
	}
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

func CheckoutBranch(name string) {
	branchPath := filepath.Join(headsDir, name)
	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		fmt.Println("Branch does not exist:", name)
		return
	}
	err := ioutil.WriteFile(headFile, []byte("ref: refs/heads/"+name+"\n"), 0644)
	if err != nil {
		fmt.Println("Error updating HEAD")
		return
	}
	fmt.Println("Switched to branch", name)
}

func UpdateHEAD(ref string) {
	err := ioutil.WriteFile(headFile, []byte(ref+"\n"), 0644)
	if err != nil {
		fmt.Println("Error updating HEAD")
	}
}

// Enhanced Push: sync objects, log, refs, HEAD
func Push(remotePath string) {
	remoteObjects := filepath.Join(remotePath, objectsDir)
	remoteLog := filepath.Join(remotePath, logFile)
	remoteRefs := filepath.Join(remotePath, refsDir)
	remoteHeads := filepath.Join(remotePath, headsDir)
	remoteHEAD := filepath.Join(remotePath, headFile)

	os.MkdirAll(remoteObjects, 0755)
	os.MkdirAll(remoteRefs, 0755)
	os.MkdirAll(remoteHeads, 0755)

	// Copy objects
	files, err := ioutil.ReadDir(objectsDir)
	if err == nil {
		for _, f := range files {
			src := filepath.Join(objectsDir, f.Name())
			dst := filepath.Join(remoteObjects, f.Name())
			data, err := ioutil.ReadFile(src)
			if err == nil {
				ioutil.WriteFile(dst, data, 0644)
			}
		}
	}

	// Copy log
	localLogData, err := ioutil.ReadFile(logFile)
	if err == nil {
		ioutil.WriteFile(remoteLog, localLogData, 0644)
	}

	// Copy refs/heads
	branches, err := ioutil.ReadDir(headsDir)
	if err == nil {
		for _, b := range branches {
			src := filepath.Join(headsDir, b.Name())
			dst := filepath.Join(remoteHeads, b.Name())
			data, err := ioutil.ReadFile(src)
			if err == nil {
				ioutil.WriteFile(dst, data, 0644)
			}
		}
	}

	// Copy HEAD
	headData, err := ioutil.ReadFile(headFile)
	if err == nil {
		ioutil.WriteFile(remoteHEAD, headData, 0644)
	}

	fmt.Println("Pushed to", remotePath)
}

// Enhanced Pull: sync objects, log, refs, HEAD
func Pull(remotePath string) {
	remoteObjects := filepath.Join(remotePath, objectsDir)
	remoteLog := filepath.Join(remotePath, logFile)
	remoteRefs := filepath.Join(remotePath, refsDir)
	remoteHeads := filepath.Join(remotePath, headsDir)
	remoteHEAD := filepath.Join(remotePath, headFile)

	os.MkdirAll(objectsDir, 0755)
	os.MkdirAll(refsDir, 0755)
	os.MkdirAll(headsDir, 0755)

	// Copy objects
	files, err := ioutil.ReadDir(remoteObjects)
	if err == nil {
		for _, f := range files {
			src := filepath.Join(remoteObjects, f.Name())
			dst := filepath.Join(objectsDir, f.Name())
			data, err := ioutil.ReadFile(src)
			if err == nil {
				ioutil.WriteFile(dst, data, 0644)
			}
		}
	}

	// Copy log
	remoteLogData, err := ioutil.ReadFile(remoteLog)
	if err == nil {
		ioutil.WriteFile(logFile, remoteLogData, 0644)
	}

	// Copy refs/heads
	branches, err := ioutil.ReadDir(remoteHeads)
	if err == nil {
		for _, b := range branches {
			src := filepath.Join(remoteHeads, b.Name())
			dst := filepath.Join(headsDir, b.Name())
			data, err := ioutil.ReadFile(src)
			if err == nil {
				ioutil.WriteFile(dst, data, 0644)
			}
		}
	}

	// Copy HEAD
	headData, err := ioutil.ReadFile(remoteHEAD)
	if err == nil {
		ioutil.WriteFile(headFile, headData, 0644)
	}

	fmt.Println("Pulled from", remotePath)
}

func Clone(remotePath, targetPath string) {
	os.MkdirAll(filepath.Join(targetPath, objectsDir), 0755)
	ioutil.WriteFile(filepath.Join(targetPath, indexFile), []byte{}, 0644)
	ioutil.WriteFile(filepath.Join(targetPath, logFile), []byte{}, 0644)
	remoteObjects := filepath.Join(remotePath, objectsDir)
	files, err := ioutil.ReadDir(remoteObjects)
	if err != nil {
		fmt.Println("Error reading remote objects")
		return
	}
	for _, f := range files {
		src := filepath.Join(remoteObjects, f.Name())
		dst := filepath.Join(targetPath, objectsDir, f.Name())
		data, err := ioutil.ReadFile(src)
		if err == nil {
			ioutil.WriteFile(dst, data, 0644)
		}
	}
	remoteLog := filepath.Join(remotePath, logFile)
	logData, err := ioutil.ReadFile(remoteLog)
	if err == nil {
		ioutil.WriteFile(filepath.Join(targetPath, logFile), logData, 0644)
	}
	fmt.Println("Cloned", remotePath, "to", targetPath)
}

func Fetch(remotePath string) {
	remoteObjects := filepath.Join(remotePath, objectsDir)
	remoteLog := filepath.Join(remotePath, logFile)
	files, err := ioutil.ReadDir(remoteObjects)
	if err != nil {
		fmt.Println("Error reading remote objects")
		return
	}
	for _, f := range files {
		dst := filepath.Join(objectsDir, f.Name())
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			src := filepath.Join(remoteObjects, f.Name())
			data, err := ioutil.ReadFile(src)
			if err == nil {
				ioutil.WriteFile(dst, data, 0644)
			}
		}
	}
	fmt.Println("Fetched objects from", remotePath)
}

func Merge(remotePath string) {
	remoteLog := filepath.Join(remotePath, logFile)
	remoteLogData, err := ioutil.ReadFile(remoteLog)
	if err != nil {
		fmt.Println("Error reading remote log")
		return
	}
	localLogData, _ := ioutil.ReadFile(logFile)
	merged := string(localLogData) + string(remoteLogData)
	ioutil.WriteFile(logFile, []byte(merged), 0644)
	fmt.Println("Merged log from", remotePath)
}

func MergeToRemote(remotePath string) {
	remoteObjects := filepath.Join(remotePath, objectsDir)
	remoteLog := filepath.Join(remotePath, logFile)

	os.MkdirAll(remoteObjects, 0755)
	if _, err := os.Stat(remoteLog); os.IsNotExist(err) {
		ioutil.WriteFile(remoteLog, []byte{}, 0644)
	}

	localFiles, err := ioutil.ReadDir(objectsDir)
	if err != nil {
		fmt.Println("Error reading local objects")
		return
	}
	for _, f := range localFiles {
		src := filepath.Join(objectsDir, f.Name())
		dst := filepath.Join(remoteObjects, f.Name())
		data, err := ioutil.ReadFile(src)
		if err == nil {
			ioutil.WriteFile(dst, data, 0644)
		}
	}

	localLogData, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Println("Error reading local log")
		return
	}
	remoteLogData, _ := ioutil.ReadFile(remoteLog)
	merged := string(remoteLogData) + string(localLogData)
	ioutil.WriteFile(remoteLog, []byte(merged), 0644)

	fmt.Println("Merged local repo into remote repo at", remotePath)
}

func Tag(name, commit string) {
	os.MkdirAll(tagsDir, 0755)
	tagPath := filepath.Join(tagsDir, name)
	if _, err := os.Stat(tagPath); err == nil {
		fmt.Println("Tag already exists:", name)
		return
	}
	err := ioutil.WriteFile(tagPath, []byte(commit), 0644)
	if err != nil {
		fmt.Println("Error creating tag:", name)
		return
	}
	fmt.Println("Created tag:", name)
}

func ListTags() {
	os.MkdirAll(tagsDir, 0755)
	files, err := ioutil.ReadDir(tagsDir)
	if err != nil {
		fmt.Println("Error reading tags")
		return
	}
	fmt.Println("Tags:")
	for _, f := range files {
		fmt.Println(" ", f.Name())
	}
}

func DeleteTag(name string) {
	tagPath := filepath.Join(tagsDir, name)
	err := os.Remove(tagPath)
	if err != nil {
		fmt.Println("Error deleting tag:", name)
		return
	}
	fmt.Println("Deleted tag:", name)
}

func ShowTag(name string) {
	tagPath := filepath.Join(tagsDir, name)
	data, err := ioutil.ReadFile(tagPath)
	if err != nil {
		fmt.Println("Tag not found:", name)
		return
	}
	fmt.Printf("Tag %s -> %s\n", name, string(data))
}

func ConfigSet(key, value string) {
	os.MkdirAll(filepath.Dir(configFile), 0755)
	lines := []string{}
	if data, err := ioutil.ReadFile(configFile); err == nil {
		lines = strings.Split(string(data), "\n")
	}
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = key + "=" + value
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, key+"="+value)
	}
	ioutil.WriteFile(configFile, []byte(strings.Join(lines, "\n")), 0644)
	fmt.Printf("Set config %s=%s\n", key, value)
}

func ConfigGet(key string) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("Config not found")
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			fmt.Println(line)
			return
		}
	}
	fmt.Println("Config key not found:", key)
}
