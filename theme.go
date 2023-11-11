package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Border(inactiveTabBorder, true).
				Padding(0, 5).
				BorderForeground(highlightColor)

	activeTabStyle = inactiveTabStyle.Copy().
			Border(activeTabBorder, true)

	windowStyle = lipgloss.NewStyle().BorderForeground(highlightColor).
			Padding(0, 4).
			Align(lipgloss.Left).
			Border(lipgloss.NormalBorder()).UnsetBorderTop()

	noStyle = lipgloss.NewStyle()

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedButton = focusedStyle.Copy().Render("[ Add User ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Add User"))

	// notification
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#32cd32"))
	failStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#dc143c"))
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
