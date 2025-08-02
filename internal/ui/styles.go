package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#7E57C2")
	secondaryColor = lipgloss.Color("#B388FF")
	subtleColor    = lipgloss.Color("#666666")
	highlightColor = lipgloss.Color("#FFD700")
	errorColor     = lipgloss.Color("#FF5252")

	// Text styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			MarginBottom(1)

	contentStyle = lipgloss.NewStyle().
			MarginLeft(2).
			MarginBottom(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// Section styles
	sectionStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1).
			MarginBottom(1).
			Width(100)

	// Table styles
	tableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0)

	// Icons
	IconTitle    = "📌 "
	IconAuthor   = "👤 "
	IconTime     = "🕒 "
	IconContent  = "📝 "
	IconComments = "💬 "
	IconRefresh  = "🔄 "
	IconBack     = "⬅️ "
	IconEnter    = "⏎ "
	IconMouse    = "🐭 "
)

var (
	StatusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true)
)
