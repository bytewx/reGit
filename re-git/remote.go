package regit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Pull(remotePath string) {
	remoteObjects := filepath.Join(remotePath, objectsDir)
	remoteLog := filepath.Join(remotePath, logFile)
	files, err := ioutil.ReadDir(remoteObjects)
	if err != nil {
		fmt.Println("Error reading remote objects")
		return
	}
	for _, f := range files {
		src := filepath.Join(remoteObjects, f.Name())
		dst := filepath.Join(objectsDir, f.Name())
		data, err := ioutil.ReadFile(src)
		if err == nil {
			ioutil.WriteFile(dst, data, 0644)
		}
	}
	remoteLogData, err := ioutil.ReadFile(remoteLog)
	if err == nil {
		localLogData, _ := ioutil.ReadFile(logFile)
		merged := string(localLogData) + string(remoteLogData)
		ioutil.WriteFile(logFile, []byte(merged), 0644)
	}
	fmt.Println("Pulled from", remotePath)
}

func Push(remotePath string) {
	remoteObjects := filepath.Join(remotePath, objectsDir)
	remoteLog := filepath.Join(remotePath, logFile)
	files, err := ioutil.ReadDir(objectsDir)
	if err != nil {
		fmt.Println("Error reading local objects")
		return
	}
	for _, f := range files {
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
	fmt.Println("Pushed to", remotePath)
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
