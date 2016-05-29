package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Command struct {
	Name, Usage string
	Execute     func([]string) error
}

const (
	cmdSymbol = "~ "
)

var (
	Commands []*Command
)

// GetCommandNames returns a slice of all command names.
func GetCommandNames() []string {
	allNames := make([]string, 0, len(Commands))
	for _, cmd := range Commands {
		allNames = append(allNames, cmd.Name)
	}
	return allNames
}

// GetCommand either returns a command or and error.
func GetCommand(name string) (*Command, error) {
	for _, cmd := range Commands {
		if cmd.Name == name {
			return cmd, nil
		}
	}
	return nil, errors.New("command not found")
}

// printMessage prints out the items with the prefix [CLI].
func printMessage(items ...interface{}) {
	fmt.Print("[CLI] ")
	fmt.Println(items...)
}

// RegisterCommand registers a new command with a name, a usage string and an executable function.
func RegisterCommand(name, usage string, executeFnc func([]string) error) {
	cmd := Command{
		Name:    name,
		Usage:   usage,
		Execute: executeFnc,
	}
	Commands = append(Commands, &cmd)
}

// init loads a standard set of commands.
func init() {
	RegisterCommand("reload", "> reload", func(args []string) error {
		printMessage("reload templates, posts and pages")
		Reload()
		return nil
	})
	RegisterCommand("stop", "> stop", func(args []string) error {
		printMessage("stop the server")
		os.Exit(1)
		return nil
	})
	RegisterCommand("debug", "> debug [mode]", func(args []string) error {
		if len(args) != 1 {
			return errors.New("bad arguments")
		}
		switch args[0] {
		case "on":
			logFlags = log.Ltime | log.Lshortfile
			initLogger(os.Stdout)
			printMessage("activated debug mode")
		case "off":
			logFlags = log.Ldate | log.Ltime
			initLogger(ioutil.Discard)
			printMessage("deactivated debug mode")
		default:
			return errors.New("unknown mode:" + args[0])
		}
		return nil
	})
	RegisterCommand("help", "> help [command]", func(args []string) error {
		if len(args) == 0 {
			printMessage("help:", strings.Join(GetCommandNames(), ", "))
		} else if len(args) == 1 {
			cmd, err := GetCommand(args[0])
			if err != nil {
				return err
			}
			printMessage("help:", cmd.Usage)
		} else {
			return errors.New("bad arguments")
		}
		return nil
	})
}

func startInteractiveMode() {
	printMessage("bloggy", BloggyVersionTag)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(cmdSymbol)
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		// Collect tokens from user input
		tokens := strings.Split(line, " ")
		for i, _ := range tokens {
			tokens[i] = strings.Trim(tokens[i], " \n\t\r")
		}

		// Just continue if no command is given
		if len(tokens) < 1 {
			continue
		} else {
			name := strings.ToLower(tokens[0])
			cmd, err := GetCommand(name)
			if err != nil {
				printMessage("error:", err)
				continue
			}

			err = cmd.Execute(tokens[1:])
			if err != nil {
				printMessage("error:", err)
				continue
			}
		}
	}
}
