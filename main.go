package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	Startup State = iota
	Ready
	Writing
	Stopped
	Resume
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
		if msg.String() == "q" {
			return m, tea.Quit
		}

		switch m.state {
		case Startup:
			// handle key events for text input
			if msg.String() == "enter" {
				if count, err := strconv.Atoi(m.textInput.Value()); err == nil {
					m.startCount = uint(count)
					m.state = Ready
					m.message = "[b]egin"
					m.textInput.Reset()
				} else {
					m.message = "invalid input. please enter a numerical wordcount."
				}
			} else {
				// pass other keys to the text input
				m.textInput, cmd = m.textInput.Update(msg)
			}
		case Ready:
			if msg.String() == "b" {
				m.startTime = time.Now()
				m.state = Writing
			}
		case Writing:
			if msg.String() == "s" {
				// save ending timestamp
				m.endTime = time.Now()
				m.message = "enter ending wordcount"

				// transition state
				m.state = Stopped
			}
		case Stopped:
			if msg.String() == "enter" {
				// retrieve ending wordcount
				if count, err := strconv.Atoi(m.textInput.Value()); err == nil {
					m.endCount = uint(count)
					m.textInput.Reset()
					m.state = Resume
				} else {
					m.message = "invalid input. please enter a numerical wordcount."
				}
			} else {
				m.textInput, cmd = m.textInput.Update(msg)
			}
		case Resume:
			if msg.String() == "r" {
				m.startCount = m.endCount
				m.startTime = time.Now()
				m.state = Writing
			}
		}
	case tea.WindowSizeMsg:
		m.textInput.Width = msg.Width / 3
	}

	return m, cmd
}

// renders the model
func (m model) View() string {
	switch m.state {
	case Startup:
		return fmt.Sprintf(
			"WPM\n\n%s\n\n%s\n\n%s",
			m.message,
			m.textInput.View(),
			"[q]uit",
		)
	case Ready:
		return fmt.Sprintf(
			"WPM\n\nbeginning wordcount: %d\n\n%s",
			m.startCount,
			"[q]uit | [b]egin",
		)
	case Writing:
		return fmt.Sprintf(
			"WPM\n\nbeginning wordcount: %d at %s\n\n%s",
			m.startCount,
			m.startTime.Format(time.Kitchen),
			"[q]uit | [s]top",
		)
	case Stopped:
		return fmt.Sprintf(
			"WPM\n\nbeginning wordcount: %d at %s\n\n%s\n\n%s\n\n%s",
			m.startCount,
			m.startTime.Format(time.Kitchen),
			m.message,
			m.textInput.View(),
			"[q]uit",
		)
	case Resume:
		durationMinutes := max(m.endTime.Sub(m.startTime).Minutes(), 1)

		return fmt.Sprintf(
			"WPM\n\nbeginning wordcount: %d at %s\n\nending wordcount: %d at %s\n\nwords per minute: %d\n\n%s",
			m.startCount,
			m.startTime.Format(time.Kitchen),
			m.endCount,
			m.endTime.Format(time.Kitchen),
			uint((m.endCount-m.startCount)/uint(durationMinutes)),
			"[q]uit | [r]esume",
		)
	default:
		return "unknown state. call for help. run in circles, scream and shout.\n\n[q]uit"
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
