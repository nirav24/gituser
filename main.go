package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	gap "github.com/muesli/go-app-paths"

	"github.com/nirav24/gituser/user"
)

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

const appName = "gituser"

func main() {
	scope := gap.NewScope(gap.User, appName)
	dirs, err := scope.DataDirs()
	if err != nil {
		fmt.Println("Error creating user store: ", err)
		os.Exit(1)
	}

	dir := dirs[0]
	fileName := "/data.json"
	if _, err := os.Stat(dir + fileName); errors.Is(err, os.ErrNotExist) {
		err := createFile(dir, fileName)
		if err != nil {
			fmt.Printf("Error creating datafile at %s%s, error %+v\n", dir, fileName, err)
			os.Exit(1)
		}
	}

	store, err := user.NewStore(dir + fileName)
	if err != nil {
		fmt.Println("Error creating user store: ", err)
		os.Exit(1)
	}
	defer func() {
		if err := store.Close(); err != nil {
			fmt.Println("Error updating users: ", err)
			os.Exit(1)
		}
	}()

	m := initialModel(store, localConfig, listView)
	program := tea.NewProgram(m)

	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func createFile(dir, name string) error {
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(dir+name, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	if _, err := file.WriteString("[]"); err != nil {
		return err
	}
	return file.Close()
}
