package main

import (
	contentmaker "contentMaker"
	"contentprocessor"
	"fmt"
)

func main() {
	go contentmaker.PopulateStorage(3, 6, 500)

	done, ticker := contentprocessor.WaitAndProces(300)
	if <-done {
		ticker.Stop()
	}
	fmt.Println("Processed and saved")
}
