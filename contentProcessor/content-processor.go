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

	for i, file := range *itemsToProcess {
		filePath := fmt.Sprintf("%v/%v", dirName, file.Name())
		content, err := os.ReadFile(filePath)
		check(err)

		if i == len(*itemsToProcess)-1 {
			concatenated += string(content)
		} else {
			concatenated += string(content) + "\n"
		}
	}

	return &concatenated
}

func ReadStorage(itemsToOmit *map[string]string) (bool, string) {
	list, err := os.ReadDir("./storage")
	check(err)

	var concatenated string

	for _, subStorage := range list {
		itemName := subStorage.Name()
		if itemName == "done.txt" {
			return true, concatenated
		}
		if len((*itemsToOmit)[itemName]) != 0 {
			continue
		}

		currentPath := fmt.Sprintf("./storage/%v", itemName)
		currentData := readSubStorage(currentPath)

		if len(*currentData) == 0 {
			log.Println("Reader is waiting...")
			continue
		}
		log.Println("Reader is reading...")
		concatenated += *currentData + "\n"
		(*itemsToOmit)[itemName] = itemName
	}

	return false, strings.Trim(concatenated, "\n")
}

func WaitAndProces(checkIntervalMs int) (chan bool, *time.Ticker) {
	ticker := time.NewTicker(time.Duration(checkIntervalMs) * time.Millisecond)
	done := make(chan bool)

	start := time.Now()
	go func() {
		itemsToOmit := make(map[string]string)
		var finalData string

		for range ticker.C {
			isFinished, concatenated := ReadStorage(&itemsToOmit)
			finalData += concatenated + "\n"

			if isFinished {
				trimmed := strings.Trim(finalData, "\n")
				os.WriteFile("./result.txt", []byte(trimmed), 0644)
				log.Printf("Waited for the final result during: %v", time.Since(start))
				done <- true
				return
			}
		}
	}()

	return done, ticker
}
