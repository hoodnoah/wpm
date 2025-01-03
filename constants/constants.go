package constants

import "github.com/charmbracelet/lipgloss"

const (
	Header = `
 _    __________  ___
| |  | | ___ \  \/  |
| |  | | |_/ / .  . |
| |/\| |  __/| |\/| |
\  /\  / |   | |  | |
 \/  \/\_|   \_|  |_/
                    `
)

var (
	TimeStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#33FF57"))
	WordCountStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	WPMStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef7c8e"))
)
