package main

import (
	// stdlib
	"fmt"
	"os"
	"strconv"
	"time"

	// external
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// internal
	"github.com/hoodnoah/wpm/m/v2/constants"
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
	spinner    spinner.Model
	message    string
	startCount uint
	endCount   uint
	startTime  time.Time
	endTime    time.Time
}

// initializes the model to its beginning state
func initialModel() model {
	// textInput
	ti := textinput.New()
	ti.Placeholder = "enter wordcount"
	ti.Focus()
	ti.CharLimit = 7
	ti.Width = 20
	ti.TextStyle = constants.WordCountStyle
	ti.PlaceholderStyle = constants.WordCountStyle.Italic(true).Faint(true)

	// spinner
	s := spinner.New()
	s.Spinner = spinner.Jump
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		state:     Startup,
		textInput: ti,
		spinner:   s,
		message:   "enter beginning wordcount",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
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

				cmd = textinput.Blink
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
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	default:
		switch m.state {
		case Writing:
			cmd = m.spinner.Tick
		case Startup, Stopped:
			m.textInput, cmd = m.textInput.Update(msg)
		}
	}

	return m, cmd
}

// renders the model
func (m model) View() string {
	switch m.state {
	case Startup:
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s\n\n%s",
			constants.Header,
			m.message,
			m.textInput.View(),
			"[q]uit",
		)
	case Ready:
		return fmt.Sprintf(
			"%s\n\nbeginning wordcount: %s\n\n%s",
			constants.Header,
			constants.WordCountStyle.Render(fmt.Sprintf("%d", m.startCount)),
			"[q]uit | [b]egin",
		)
	case Writing:
		return fmt.Sprintf(
			"%s\n\nbeginning wordcount: %s at %s\n\n%s writing...\n\n%s",
			constants.Header,
			constants.WordCountStyle.Render(fmt.Sprintf("%d", m.startCount)),
			constants.TimeStyle.Render(m.startTime.Format(time.Kitchen)),
			m.spinner.View(),
			"[q]uit | [s]top",
		)
	case Stopped:
		return fmt.Sprintf(
			"%s\n\nbeginning wordcount: %s at %s\n\n%s\n\n%s\n\n%s",
			constants.Header,
			constants.WordCountStyle.Render(fmt.Sprintf("%d", m.startCount)),
			constants.TimeStyle.Render(m.startTime.Format(time.Kitchen)),
			m.message,
			m.textInput.View(),
			"[q]uit",
		)
	case Resume:
		durationMinutes := max(m.endTime.Sub(m.startTime).Minutes(), 1)

		return fmt.Sprintf(
			"%s\n\nbeginning wordcount: %s at %s\n\nending wordcount: %s at %s\n\nwords per minute: %s\n\n%s",
			constants.Header,
			constants.WordCountStyle.Render(fmt.Sprintf("%d", m.startCount)),
			constants.TimeStyle.Render(m.startTime.Format(time.Kitchen)),
			constants.WordCountStyle.Render(fmt.Sprintf("%d", m.endCount)),
			constants.TimeStyle.Render(m.endTime.Format(time.Kitchen)),
			constants.WPMStyle.Render(fmt.Sprintf("%d", uint((m.endCount-m.startCount)/uint(durationMinutes)))),
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
