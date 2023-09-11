package main

import (
	"bufio"
	contentmaker "contentMaker"
	"contentprocessor"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Params struct {
	dirsCount      int
	filesCount     int
	minimalDelayMs int
}

func getStartParams(paramsNames []string) *Params {
	var params Params
	reader := bufio.NewReader(os.Stdin)

	for _, name := range paramsNames {
		var label string
		switch name {
		case "dirsCount":
			label = "Please enter desired sub-storages count: "
		case "filesCount":
			label = "Please enter desired files count in every sub-storage: "
		case "minimalDelayMs":
			label = "Please enter minimal delay in ms for the timer (ticker): "
		}

		fmt.Fprint(os.Stderr, label)

		p, _ := reader.ReadString('\n')
		v, err := strconv.ParseInt(strings.TrimSpace(p), 0, 0)

		if err != nil {
			panic(err)
		}

		switch name {
		case "dirsCount":
			params.dirsCount = int(v)
		case "filesCount":
			params.filesCount = int(v)
		case "minimalDelayMs":
			params.minimalDelayMs = int(v)
		}
	}

	return &params
}

func main() {
	paramsNames := []string{"dirsCount", "filesCount", "minimalDelayMs"}
	params := getStartParams(paramsNames)

	go contentmaker.PopulateStorage(params.dirsCount, params.filesCount, params.minimalDelayMs)

	done, ticker := contentprocessor.WaitAndProces(300)
	if <-done {
		ticker.Stop()
	}
	fmt.Println("Processed and saved")
}
