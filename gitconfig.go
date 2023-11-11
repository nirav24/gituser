package main

import (
	"fmt"
	"github.com/nirav24/gituser/user"
	"os/exec"
	"reflect"
)

func setConfig(u user.User, cMode configMode) error {

	args := []string{"config"}

	if cMode == globalConfig {
		args = append(args, "--global")
	} else {
		// Check if current directory is git repo
		cmd := exec.Command("git", "status")
		err := cmd.Run()
		if err != nil {
			return err
		}

		if cmd.ProcessState.ExitCode() != 0 {
			return fmt.Errorf("not a git repository")
		}
	}

	rv := reflect.ValueOf(u)
	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		cmdTag, ok := t.Field(i).Tag.Lookup("cmd")
		if !ok || cmdTag == "" {
			continue // if tag is empty
		}

		newArgs := make([]string, 0)
		newArgs = append(newArgs, args...)

		if rv.Field(i).Kind() == reflect.String {
			if rv.Field(i).String() != "" {
				newArgs = append(newArgs, cmdTag, rv.Field(i).String())
			} else if isValueSet(cmdTag, cMode) {
				newArgs = append(newArgs, "--unset", cmdTag)
			}
		} else if rv.Field(i).Kind() == reflect.Bool {
			if rv.Field(i).Bool() {
				newArgs = append(newArgs, cmdTag, "true")
			} else if isValueSet(cmdTag, cMode) {
				newArgs = append(newArgs, "--unset", cmdTag)
			}
		}
		// if we did not add any new arg to newArgs,
		// then we do not need to run any command for given cmd tag
		if len(newArgs) == len(args) {
			continue
		}

		cmd := exec.Command("git", newArgs...)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to update %s", cmdTag)
		}

	}
	return nil
}

func isValueSet(key string, cMode configMode) bool {
	args := []string{"config"}

	if cMode == globalConfig {
		args = append(args, "--global")
	} else {
		args = append(args, "--local")
	}

	args = append(args, "--get", key)
	cmd := exec.Command("git", args...)
	err := cmd.Run()
	if err != nil && cmd.ProcessState.ExitCode() == 1 {
		return false
	}

	return true
}
