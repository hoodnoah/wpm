package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// input beginning word count
// [b] begin
// [s] stop
// input end word count
// display rounded wpm

type State int

const (
	Startup State = iota
	Ready
	Writing
	Stopped
)

type model struct {
	state      State
	textInput  textinput.Model
	message    string
	startCount uint
	endCount   uint
	startTime  time.Time
	endTime    time.Time
}

// initializes the model to its beginning state
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "enter wordcount"
	ti.Focus()
	ti.CharLimit = 7
	ti.Width = 20

	return model{
		state:     Startup,
		textInput: ti,
		message:   "enter beginning wordcount",
	}
}

// init returns a command used for initial I/O.
// I have none, so we can return nil.
func (m model) Init() tea.Cmd {
	// return nil, meaning "no I/O atm"
	return nil
}

// updates the underlying model based on action
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// is it a keypress?
	case tea.KeyMsg:
		switch m.state {
		case Startup:
			// handle key events for text input
			if msg.String() == "enter" {
				if count, err := strconv.Atoi(m.textInput.Value()); err == nil {
					m.StartCount = count
					m.State = Ready
					m.message = "[b]egin"
					m.textInput.Reset()
				} else {
					m.message = "invalid input. please enter a numerical wordcount."
				}
			} else {
				// pass other keys to the text input
				m.textInput, cmd = m.textInput.Update(msg)
			}
		}
	case tea.WindowSizeMsg:
		m.textInput.Width = msg.Width / 3
	}

	return m, cmd
}

// renders the model
func (m model) View() string {
	var inputView string
	if m.state == Startup || m.state == Stopped {
		inputView = m.textInput.View()
	}

	return fmt.Sprintf(
		"WPM\n\n%s\n\n%s\n\n%s\n",
		m.message,
		inputView,
		"[q]uit",
	)
}

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
