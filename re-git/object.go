package regit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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

func RemoveObject(oid string) {
	objPath := filepath.Join(objectsDir, oid)
	err := os.Remove(objPath)
	if err != nil {
		fmt.Println("Error removing object:", oid)
		return
	}
	fmt.Println("Removed object:", oid)
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
