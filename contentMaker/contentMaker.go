package contentmaker

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func AddNewContentPart(dirNum, filesCount, delayMs int) {
	ticker := time.NewTicker(time.Duration(delayMs) * time.Millisecond)
	done := make(chan bool)

	dirName := fmt.Sprintf("./storage/sub-storage-%v", dirNum)
	errRem := os.RemoveAll(dirName)
	check(errRem)

	err := os.Mkdir(dirName, 0755)
	check(err)

	start := time.Now()
	go func() {
		counter := 0
		for {
			select {
			case <-done:
				finalFile, err := os.Create(fmt.Sprintf("%v/done.txt", dirName))
				check(err)

				defer finalFile.Close()
				return

			case t := <-ticker.C:
				counter += 1
				stubContent := []byte(t.String())
				destination := fmt.Sprintf("%v/file-%v.txt", dirName, counter)

				err := os.WriteFile(destination, stubContent, 0644)
				check(err)
			}
		}
	}()

	time.Sleep(time.Duration(filesCount*delayMs+100) * time.Millisecond)
	ticker.Stop()
	done <- true

	log.Printf("Initial content making for directory num.%v is finished.\nProcess took: %v\n\n", dirNum, time.Since(start))
}

func PopulateStorage(dirsCount, filesCount, delayMs int) {
	errRem := os.RemoveAll("./storage")
	check(errRem)

	err := os.Mkdir("./storage", 0755)
	check(err)

	var wg sync.WaitGroup

	for i := 0; i < dirsCount; i += 1 {
		wg.Add(1)
		go func(dirNum int) {
			defer wg.Done()
			AddNewContentPart(dirNum, filesCount, delayMs)
		}(i)
	}
	wg.Wait()
	finalFile, err := os.Create("./storage/done.txt")
	check(err)

	defer finalFile.Close()
}
