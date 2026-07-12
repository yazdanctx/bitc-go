package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	fileSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#04B575")).
				Bold(true)

	fileDefaultStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA"))

	checkboxOnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	checkboxOffStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#626262"))

	startButtonStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 3)

	startButtonInactiveStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#626262")).
					Padding(0, 3)

	progressBarStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4"))

	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4")).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#626262"))

	tableCellStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	summaryStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B"))

	goodSavingsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#04B575"))

	badSavingsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B"))
)

func FormatSize(bytes int64) string {
	switch {
	case bytes >= 1<<20:
		return formatFloat(float64(bytes)/(1<<20)) + " MB"
	case bytes >= 1<<10:
		return formatFloat(float64(bytes)/(1<<10)) + " KB"
	default:
		return formatFloat(float64(bytes)) + " B"
	}
}

func formatFloat(f float64) string {
	s := "%.1f"
	if f >= 100 {
		s = "%.0f"
	}
	return fmt.Sprintf(s, f)
}

func FormatSavings(pct float64) string {
	if pct >= 0 {
		return goodSavingsStyle.Render(fmt.Sprintf("-%.1f%%", pct))
	}
	return badSavingsStyle.Render(fmt.Sprintf("+%.1f%%", -pct))
}
