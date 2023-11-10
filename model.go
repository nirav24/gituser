package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nirav24/gituser/user"
	"strings"
)

type model struct {
	// main model state
	mode      view
	activeTab int
	Err       error
	cMode     configMode

	// Show list of users
	listModel list.Model

	// Add User
	inputs     []textinput.Model
	focusIndex int
}

func initialModel(cMode configMode, currentView view) *model {
	items := make([]list.Item, len(users))
	for i, u := range users {
		items[i] = u
	}
	listModel := list.New(items, itemDelegate{}, 60, 30)
	listModel.SetShowStatusBar(false)
	listModel.SetShowPagination(false)
	listModel.AdditionalShortHelpKeys = userListKeys
	listModel.AdditionalFullHelpKeys = userListKeys

	m := model{
		inputs:    make([]textinput.Model, 3),
		listModel: listModel,
		mode:      listView,
		cMode:     localConfig,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Nickname"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Email"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Signing Key (Optional): "
		}

		m.inputs[i] = t
	}

	return &m
}

func (m *model) saveNewUser() tea.Cmd {
	u := user.User{
		Username:   m.inputs[0].Value(),
		Email:      m.inputs[1].Value(),
		SigningKey: m.inputs[2].Value(),
		ShouldSign: false,
	}
	users = append(users, u)
	return m.listModel.InsertItem(len(m.listModel.Items()), u)
}

func (m *model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "right", "tab", "left", "shift+tab":
			if m.mode == listView {
				m.mode = addView
			} else {
				m.mode = listView
			}
			return m, nil
		case "ctrl+g":
			if m.cMode == localConfig {
				m.cMode = globalConfig
			} else {
				m.cMode = localConfig
			}
			return m, nil
		}
	}

	if m.mode == listView {
		var cmd tea.Cmd
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				currentUser := m.listModel.SelectedItem().(user.User)
				return m, tea.Sequence(m.listModel.NewStatusMessage(fmt.Sprintf("User: %s is set", currentUser.Username)), tea.Quit)
			case "backspace":
				// only invoke list update if there is any item
				if len(m.listModel.Items()) > 0 {
					removedUser := m.listModel.SelectedItem().(user.User)
					m.listModel.RemoveItem(m.listModel.Index())
					return m, m.listModel.NewStatusMessage(fmt.Sprintf("User: %s is removed", removedUser.Username))
				}
			}
		}
		m.listModel, cmd = m.listModel.Update(msg)

		return m, cmd
	}

	return m.updateInputs(msg)
}

func (m *model) View() string {
	doc := strings.Builder{}

	var renderedTabs []string
	tabs := []string{"Select User", "Add User"}
	for i, t := range tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(tabs)-1, i == int(m.mode)
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Width(30-windowStyle.GetHorizontalFrameSize()).Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	if m.mode == listView {
		// Set title here to show config mode dynamically
		m.listModel.Title = fmt.Sprintf("List of Users (%s config)", m.cMode)
		doc.WriteString(windowStyle.Width(60 - windowStyle.GetHorizontalFrameSize()).
			Render(m.listModel.View()))
	} else {
		doc.WriteString(windowStyle.Width(60 - windowStyle.GetHorizontalFrameSize()).
			Render(m.viewAddUserForm()))
	}

	return docStyle.Render(doc.String())
}

// Add user related functions
func (m *model) viewAddUserForm() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}

	fmt.Fprintf(&b, "\n\n\n%s\n", *button)

	return b.String() + "\n\n" + help.New().ShortHelpView(addUserKeys())
}

func (m *model) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			m.resetAddForm()

			return m, tea.Batch(cmds...)
		// Set focus to next input
		case "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				cmd := m.saveNewUser()
				// reset to view mode
				m.mode = listView
				m.resetAddForm()

				return m, cmd
			}

			// Cycle indexes
			if s == "up" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
				} else {
					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

// helper function
func (m *model) resetAddForm() {
	m.focusIndex = 0
	for i := range m.inputs {
		m.inputs[i].Reset()
		if i == m.focusIndex {
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = blurredStyle
			m.inputs[i].TextStyle = blurredStyle
		}
	}

}
