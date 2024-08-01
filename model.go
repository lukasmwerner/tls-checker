package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	xMark = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).SetString("X")
)

type successMsg string
type failureMsg string

type model struct {
	domain string

	processingTLS10 bool
	statusTLS10 string
	processingTLS11 bool
	statusTLS11 string
	processingTLS12 bool
	statusTLS12 string
	processingTLS13 bool
	statusTLS13 string

	spinner spinner.Model
}
func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}
	case successMsg:
		switch msg {
		case "TLS1.0":
			m.statusTLS10 = ""
			m.processingTLS10 = false
		case "TLS1.1":
			m.statusTLS11 = ""
			m.processingTLS11 = false
		case "TLS1.2":
			m.statusTLS12 = ""
			m.processingTLS12 = false
		case "TLS1.3":
			m.statusTLS12 = ""
			m.processingTLS13 = false
		}
		return m, nil
	case failureMsg:
		switch msg[len(msg)-1] {
		case '0':
			m.statusTLS10 = string(msg)[:len(msg)-7]
			m.processingTLS10 = false
		case '1':
			m.statusTLS11 = string(msg)[:len(msg)-7]
			m.processingTLS11 = false
		case '2':
			m.statusTLS12 = string(msg)[:len(msg)-7]
			m.processingTLS12 = false
		case '3':
			m.statusTLS13 = string(msg)[:len(msg)-7]
			m.processingTLS13 = false
		}
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	render := " Checking: " +  m.domain + "\n"
	if m.processingTLS10 {
		render += fmt.Sprintf(" %s TLS1.0\n", m.spinner.View())
	} else {
		if m.statusTLS10 == "" {
			render += fmt.Sprintf("  %s  TLS1.0\n", checkMark.String())
		} else {
			render += fmt.Sprintf("  %s  TLS1.0: %s\n", xMark.String(), m.statusTLS10)
		}
	}


	if m.processingTLS11 {
		render += fmt.Sprintf(" %s TLS1.1\n", m.spinner.View())
	} else {
		if m.statusTLS11 == "" {
			render += fmt.Sprintf("  %s  TLS1.1\n", checkMark.String())
		} else {
			render += fmt.Sprintf("  %s  TLS1.1: %s \n", xMark.String(), m.statusTLS11)
		}
	}


	if m.processingTLS12 {
		render += fmt.Sprintf(" %s TLS1.2\n", m.spinner.View())
	} else {
		if m.statusTLS12 == "" {
			render += fmt.Sprintf("  %s  TLS1.2\n", checkMark.String())
		} else {
			render += fmt.Sprintf("  %s  TLS1.2: %s \n", xMark.String(), m.statusTLS12)
		}
	}


	if m.processingTLS13 {
		render += fmt.Sprintf(" %s TLS1.3\n", m.spinner.View())
	} else {
		if m.statusTLS13 == "" {
			render += fmt.Sprintf("  %s  TLS1.3\n", checkMark.String())
		} else {
			render += fmt.Sprintf("  %s  TLS1.3: %s \n", xMark.String(), m.statusTLS13)
		}
	}
			
	return render
}

func initalModel(hostname string) model {
	m := model{
		domain:          hostname,
		processingTLS10: true,
		statusTLS10:     "",
		processingTLS11: true,
		statusTLS11:     "",
		processingTLS12: true,
		statusTLS12:     "",
		processingTLS13: true,
		statusTLS13:     "",
		spinner:         spinner.New(spinner.WithSpinner(spinner.Points)),
	}
	return m
}
