package contentprocessor

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkIfReady(dirName string) *[]fs.DirEntry {
	list, err := os.ReadDir(dirName)
	check(err)

	var itemsToProcess []fs.DirEntry

	ready := false
	for _, item := range list {
		if item.Name() == "done.txt" {
			ready = true
		} else {
			itemsToProcess = append(itemsToProcess, item)
		}
	}
	if !ready {
		return nil
	}
	return &itemsToProcess
}

func readSubStorage(dirName string) *string {
	var concatenated string
	var itemsToProcess []fs.DirEntry = *checkIfReady(dirName)

	if itemsToProcess == nil {
		return &concatenated
	}

	for i, file := range itemsToProcess {
		filePath := fmt.Sprintf("%v/%v", dirName, file.Name())
		content, err := os.ReadFile(filePath)
		check(err)

		if i == len(itemsToProcess)-1 {
			concatenated += string(content)
		} else {
			concatenated += string(content) + "\n"
		}
	}

	return &concatenated
}

func ReadStorage() string {
	subStorages := checkIfReady("./storage")
	var finalData string

	if subStorages == nil {
		log.Println("Reader is waiting...")
		return finalData
	}

	for i, subSt := range *subStorages {
		currentPath := fmt.Sprintf("./storage/%v", subSt.Name())
		currentData := readSubStorage(currentPath)

		if len(*currentData) == 0 {
			continue
		}
		if i == len(*subStorages)-1 {
			finalData += *currentData
		} else {
			finalData += *currentData + "\n"
		}
	}

	return finalData
}

func WaitAndProces(checkIntervalMs int) (chan bool, *time.Ticker) {
	ticker := time.NewTicker(time.Duration(checkIntervalMs) * time.Millisecond)
	done := make(chan bool)

	start := time.Now()
	go func() {
		for range ticker.C {
			concatenated := ReadStorage()
			if len(concatenated) != 0 {
				os.WriteFile("./result.txt", []byte(concatenated), 0644)
				log.Printf("Waited for the final result during: %v", time.Since(start))
				done <- true
				return
			}
		}
	}()

	return done, ticker
}
