package contentprocessor

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
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
	var itemsToProcess *[]fs.DirEntry = checkIfReady(dirName)

	if itemsToProcess == nil {
		return &concatenated
	}

	for _, file := range *itemsToProcess {
		filePath := fmt.Sprintf("%v/%v", dirName, file.Name())
		content, err := os.ReadFile(filePath)
		check(err)

		concatenated += string(content) + "\n"
	}

	return &concatenated
}

func ReadStorage(itemsToOmit *map[string]string) (bool, string) {
	list, err := os.ReadDir("./storage")
	check(err)

	var itemsToProcess []fs.DirEntry
	recordingCompleted := false

	for _, subStorage := range list {
		itemName := subStorage.Name()
		if len((*itemsToOmit)[itemName]) != 0 {
			continue
		}
		if itemName == "done.txt" {
			recordingCompleted = true
		} else {
			itemsToProcess = append(itemsToProcess, subStorage)
		}
	}

	var concatenated string
	for _, subStItem := range itemsToProcess {
		itemName := subStItem.Name()
		currentPath := fmt.Sprintf("./storage/%v", itemName)
		currentData := readSubStorage(currentPath)

		if len(*currentData) == 0 {
			continue
		}
		log.Printf("Reader is reading sub-storage %v...", itemName)
		concatenated += *currentData
		(*itemsToOmit)[itemName] = itemName
	}

	if !recordingCompleted {
		log.Println("Reader is waiting...")
		return false, concatenated
	}

	return true, concatenated
}

func WaitAndProces(checkIntervalMs int) (chan bool, *time.Ticker) {
	ticker := time.NewTicker(time.Duration(checkIntervalMs) * time.Millisecond)
	done := make(chan bool)

	start := time.Now()
	go func(ch chan bool) {
		itemsToOmit := make(map[string]string)
		var finalData string

		for range ticker.C {
			isFinished, concatenated := ReadStorage(&itemsToOmit)
			finalData += concatenated

			if isFinished {
				trimmed := strings.Trim(finalData, "\n")
				os.WriteFile("./result.txt", []byte(trimmed), 0644)
				log.Printf("Waited for the final result during: %v", time.Since(start))
				ch <- true
				return
			}
		}
	}(done)

	return done, ticker
}
