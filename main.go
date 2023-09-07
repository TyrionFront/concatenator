package main

import (
	contentmaker "contentMaker"
	"contentprocessor"
	"fmt"
)

func main() {
	go contentmaker.PopulateStorage(3, 4, 1000)

	done, ticker := contentprocessor.WaitAndProces(100)
	if <-done {
		ticker.Stop()
	}
	fmt.Println("Processed and saved")
}
