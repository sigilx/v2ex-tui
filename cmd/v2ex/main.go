package main

import (
	"fmt"
	"time"

	"v2ex-tui/internal/ui"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

type page int

const (
	homeView page = iota
	detailView
)

type model struct {
	currentPage   page
	homePage      *ui.HomePage
	detailPage    *ui.DetailPage
	mouseEnabled  bool
	statusMessage string
	timer         *time.Timer
}

type statusMessageTimeout struct{}

func initialModel() model {
	return model{
		currentPage:  homeView,
		homePage:     ui.NewHomePage(),
		detailPage:   ui.NewDetailPage(),
		mouseEnabled: true, // Enable mouse by default
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.homePage.Init(), tea.EnableMouseCellMotion)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "backspace", " ", "left", "h":
			if m.currentPage == detailView {
				m.currentPage = homeView
				m.detailPage.Reset()
				return m, nil
			}
		case "enter", "right", "l":
			if m.currentPage == homeView {
				if topic := m.homePage.GetSelectedTopic(); topic != nil {
					m.currentPage = detailView
					return m, m.detailPage.LoadTopic(*topic)
				}
			}
		case "m":
			m.mouseEnabled = !m.mouseEnabled
			if m.mouseEnabled {
				return m, tea.EnableMouseCellMotion
			}
			return m, tea.DisableMouse
		case "f":
			if m.currentPage == homeView {
				if topic := m.homePage.GetSelectedTopic(); topic != nil {
					clipboard.WriteAll(topic.URL)
					m.statusMessage = "Copied to clipboard: " + topic.URL
					if m.timer != nil {
						m.timer.Stop()
					}
					m.timer = time.NewTimer(2 * time.Second)
					return m, func() tea.Msg {
						<-m.timer.C
						return statusMessageTimeout{}
					}
				}
			}
		}
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			if m.currentPage == homeView {
				if topic := m.homePage.GetSelectedTopic(); topic != nil {
					m.currentPage = detailView
					return m, m.detailPage.LoadTopic(*topic)
				}
			}
		}
	case statusMessageTimeout:
		m.statusMessage = ""
		return m, nil
	}

	var cmd tea.Cmd
	switch m.currentPage {
	case homeView:
		var homePage *ui.HomePage
		homePage, cmd = m.homePage.Update(msg)
		m.homePage = homePage
	case detailView:
		var detailPage *ui.DetailPage
		detailPage, cmd = m.detailPage.Update(msg)
		m.detailPage = detailPage
	}
	return m, cmd
}

func (m model) View() string {
	var s string
	switch m.currentPage {
	case homeView:
		s = m.homePage.View()
	case detailView:
		s = m.detailPage.View()
	default:
		s = "Unknown view"
	}
	if m.statusMessage != "" {
		return s + "\n" + ui.StatusMessageStyle.Render(m.statusMessage)
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		return
	}
}
