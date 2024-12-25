package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// input beginning word count
// [b] begin
// [s] stop
// input end word count
// display rounded wpm

type WordStats struct {
	StartCount int
	EndCount   int
	StartTime  time.Time
	EndTime    time.Time
}

type Command int

const (
	Begin = iota
	Stop
)

// parses user command input into the given command enum member
// b -> begin
// s -> stop
func parseUserCommand(userInput string) (Command, error) {
	trimmedInput := strings.ToLower(userInput)

	if len(trimmedInput) > 1 {
		return 0, fmt.Errorf("expected a single character input, received %d chars", len(userInput))
	}

	switch trimmedInput {
	case "b":
		return Begin, nil
	case "s":
		return Stop, nil
	default:
		return 0, fmt.Errorf("unidentified command input: %s", trimmedInput)
	}
}

func parseInt(wordcountString string) (int, error) {
	wordCount, err := strconv.Atoi(wordcountString)
	if err != nil {
		return 0, err
	} else {
		return wordCount, nil
	}
}

func loop() {

}

func main() {
	var stats WordStats

	var wordCountString string
	var commandString string

	fmt.Print("beginning wordcount: ")
	fmt.Scan(&wordCountString)
	wordCount, err := parseInt(wordCountString)
	if err != nil {
		fmt.Printf("%s", err)
	}
	stats.StartCount = wordCount

	fmt.Print("[b]egin ")
	fmt.Scan(&commandString)
	cmd, err := parseUserCommand(commandString)
	if err != nil {
		fmt.Printf("%v", err)
	} else if cmd != Begin {
		fmt.Println("expected b to begin")
	}
	stats.StartTime = time.Now()

	fmt.Println("Started at ", stats.StartTime.Format(time.Kitchen))

	fmt.Print("[s]top ")
	fmt.Scan(&commandString)
	cmd, err = parseUserCommand(commandString)
	if err != nil {
		fmt.Printf("%v", err)
	} else if cmd != Stop {
		fmt.Println("expected s to stop")
	}
	stats.EndTime = time.Now()
	fmt.Println("Ended at ", stats.EndTime.Format(time.Kitchen))

	fmt.Print("ending wordcount: ")
	fmt.Scan(&wordCountString)
	wordCount, err = parseInt(wordCountString)
	if err != nil {
		fmt.Printf("%s", err)
	}
	stats.EndCount = wordCount

	duration := stats.EndTime.Sub(stats.StartTime)
	durationMinutes := duration.Minutes()
	wpm := float64(stats.EndCount-stats.StartCount) / durationMinutes
	fmt.Printf("%.0f wpm\n", wpm)

}
