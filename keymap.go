package main

import (
	"github.com/charmbracelet/bubbles/key"
)

var (
	nextTab = key.NewBinding(
		key.WithKeys("tab", "right"),
		key.WithHelp("tab/→", "switch tab"),
	)

	selectUserEnterKey = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	)

	deleteUserEnterKey = key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "delete"),
	)

	configModeChange = key.NewBinding(
		key.WithKeys("ctrl+g"),
		key.WithHelp("ctrl+g", "switch git config type"),
	)

	resetForm = key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "Reset Form"),
	)
	quit = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)
)

func userListKeys() []key.Binding {
	return []key.Binding{
		nextTab,
		configModeChange,
		selectUserEnterKey,
		deleteUserEnterKey,
	}
}

func addUserKeys() []key.Binding {
	return []key.Binding{
		nextTab,
		resetForm,
		quit,
	}
}
