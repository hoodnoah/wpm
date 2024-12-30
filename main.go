package main

import (
	"fmt"
	"os"
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
	Quit
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
	case "q":
		return Quit, nil
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

// query the user for wordcount, parse into integer and store in provided struct
func getWordCount(stats *WordStats) error {
	var countString string
	if stats.StartCount != 0 {
		fmt.Print("ending wordcount: ")
	} else {
		fmt.Print("starting wordcount: ")
	}

	fmt.Scan(&countString)
	intCount, err := parseInt(countString)
	if err != nil {
		return err
	}

	if stats.StartCount != 0 {
		stats.EndCount = intCount
	} else {
		stats.StartCount = intCount
	}

	return nil
}

func loop() {

}

func main() {
	var stats WordStats

	var commandString string

	if err := getWordCount(&stats); err != nil {
		fmt.Printf("%v", err)
	}

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

	if err := getWordCount(&stats); err != nil {
		fmt.Printf("%v", err)
	}

	duration := stats.EndTime.Sub(stats.StartTime)
	durationMinutes := duration.Minutes()
	wpm := float64(stats.EndCount-stats.StartCount) / durationMinutes
	fmt.Printf("%.0f wpm\n", wpm)

	fmt.Printf("[q]uit ")
	fmt.Scan(&commandString)
	cmd, err = parseUserCommand(commandString)
	if cmd == Quit {
		os.Exit(0)
	} else {
		fmt.Printf("expected q to quit, received %s", commandString)
	}

}
