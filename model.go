package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nirav24/gituser/user"
)

type model struct {
	// main model state
	store      *user.Store
	activeView view
	Err        error
	cMode      configMode

	height int
	width  int

	// Show list of users
	listModel list.Model

	// Add User
	inputs     []textinput.Model
	focusIndex int
}

func initialModel(store *user.Store, cMode configMode, currentView view) *model {
	users := store.GetUsers()

	items := make([]list.Item, len(users))
	for i, u := range users {
		items[i] = u
	}
	listModel := list.New(items, list.NewDefaultDelegate(), 0, 0)
	listModel.SetShowStatusBar(false)
	listModel.SetShowPagination(false)
	listModel.AdditionalShortHelpKeys = userListKeys
	listModel.AdditionalFullHelpKeys = userListKeys

	m := model{
		store:      store,
		inputs:     make([]textinput.Model, 5),
		listModel:  listModel,
		activeView: currentView,
		cMode:      cMode,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Username: "
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Email: "
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Signing Key (Optional): "
		case 3:
			t.Placeholder = "Sign Commits: (Yes / No)"
		case 4:
			t.Placeholder = "Format: (openpgp, ssh)"
		}

		m.inputs[i] = t
	}

	return &m
}

func (m *model) saveNewUser() tea.Cmd {
	var signCommits *bool
	if strings.ToLower(m.inputs[3].Value()) == "yes" {
		signCommits = pointer(true)
	} else if m.inputs[3].Value() != "" {
		signCommits = pointer(false)
	}

	u := user.User{
		Username:    m.inputs[0].Value(),
		Email:       m.inputs[1].Value(),
		SigningKey:  m.inputs[2].Value(),
		SignCommits: signCommits,
		GpgFormat:   m.inputs[4].Value(),
	}
	m.Err = m.store.AddUser(u)
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
			if m.activeView == listView {
				m.activeView = addView
			} else {
				m.activeView = listView
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
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.listModel.SetWidth(msg.Width)
		m.listModel.SetHeight(msg.Height - 20)
	}

	if m.activeView == listView {
		var cmd tea.Cmd
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				// only invoke if there is any item in list
				if len(m.listModel.Items()) == 0 {
					return m, nil
				}
				currentUser := m.listModel.SelectedItem().(user.User)
				m.Err = setConfig(currentUser, m.cMode)
				if m.Err != nil {
					return m.showAndResetErr()
				}
				return m, tea.Sequence(m.successMessage(fmt.Sprintf("%s is set", currentUser.Username)))
			case "backspace":
				// only invoke if there is any item in list
				if len(m.listModel.Items()) > 0 {
					removedUser := m.listModel.SelectedItem().(user.User)
					m.listModel.RemoveItem(m.listModel.Index())
					m.store.RemoveUser(removedUser)
					return m, m.successMessage(fmt.Sprintf("User: %s is removed", removedUser.Username))
				}
			}
		}
		m.listModel, cmd = m.listModel.Update(msg)
		if m.Err != nil {
			return m.showAndResetErr()
		}
		return m, cmd
	}

	// if view is addView, then call updateInputs
	return m.updateInputs(msg)
}

func (m *model) View() string {
	doc := strings.Builder{}

	var renderedTabs []string
	tabs := []string{"Select User", "Add User"}
	width := m.width

	if width%2 == 1 {
		width = width - 1
	}

	for i, t := range tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(tabs)-1, i == int(m.activeView)
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
		style = style.Border(border).
			Width((width - style.GetHorizontalFrameSize()) / 2)

		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	style := windowStyle.Copy().Width(width - windowStyle.GetHorizontalFrameSize())

	if m.activeView == listView {
		// Set title here to show config activeView dynamically
		m.listModel.Title = fmt.Sprintf("List of Users (%s config)", m.cMode)
		doc.WriteString(style.Render(m.listModel.View()))
	} else {
		doc.WriteString(style.Render(m.viewAddUserForm()))
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
				// reset to view activeView
				m.activeView = listView
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

func (m *model) successMessage(message string) tea.Cmd {
	return m.listModel.NewStatusMessage(successStyle.Render(message))
}

func (m *model) failMessage(message string) tea.Cmd {
	return m.listModel.NewStatusMessage(failStyle.Render(message))
}

func (m *model) showAndResetErr() (tea.Model, tea.Cmd) {
	if m.Err == nil {
		return m, nil
	}

	msg := m.Err.Error()
	m.Err = nil // reset error
	return m, m.failMessage(msg)
}

func pointer[T any](value T) *T {
	return &value
}
