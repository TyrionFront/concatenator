package contentmaker

import (
	"fmt"
	"log"
	"math/rand"
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
	ticker := time.NewTicker(time.Duration(delayMs*(rand.Intn(dirNum)+1)) * time.Millisecond)
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
				err := os.WriteFile(fmt.Sprintf("%v/done.txt", dirName), []byte(dirName), 0644)
				check(err)
				return

			case t := <-ticker.C:
				counter += 1
				stubContent := []byte(fmt.Sprintf("%v - substorage N %v", t.String(), dirNum))
				destination := fmt.Sprintf("%v/file-%v.txt", dirName, counter)

				err := os.WriteFile(destination, stubContent, 0644)
				check(err)
			}
		}
	}()

	time.Sleep(time.Duration(filesCount*delayMs) * time.Millisecond)
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

	for i := 1; i <= dirsCount; i += 1 {
		wg.Add(1)
		go func(dirNum int) {
			defer wg.Done()
			AddNewContentPart(dirNum, filesCount, delayMs)
		}(i)
	}
	wg.Wait()
	time.Sleep(500 * time.Millisecond)
	err = os.WriteFile("./storage/done.txt", []byte("done"), 0644)
	check(err)
}
