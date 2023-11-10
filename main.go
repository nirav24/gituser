package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nirav24/gituser/user"
	"os"
)

var users = []user.User{
	{
		Username:   "n1",
		Email:      "email1",
		SigningKey: "",
	},
	{
		Username:   "n2",
		Email:      "email2",
		SigningKey: "",
	},
}

type configMode string

const (
	globalConfig configMode = "Global"
	localConfig  configMode = "Local"
)

type view int

const (
	listView view = iota
	addView
)

func main() {
	if _, err := tea.NewProgram(initialModel(localConfig, listView)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
