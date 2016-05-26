package main

import (
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
	cmdSymbol = "$ "
)

var (
	commands = make(map[string]Command)
)

// RegisterCommand registers a new command with a name, a usage string and an executable function.
func RegisterCommand(name, usage string, executeFnc func([]string) error) {
	cmd := Command{
		Name:    name,
		Usage:   usage,
		Execute: executeFnc,
	}
	commands[name] = cmd
}

// init loads a standard set of commands.
func init() {
	RegisterCommand("reload", "> reload", func(args []string) error {
		if err := LoadTemplates(); err != nil {
			return err
		}
		if err := LoadPosts(); err != nil {
			return err
		}
		if err := LoadPages(); err != nil {
			return err
		}
		fmt.Println("reload templates, posts and pages")
		return nil
	})
	RegisterCommand("stop", "> stop", func(args []string) error {
		fmt.Println("stop the server")
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
			Info.Println("activated debug mode")
		case "off":
			logFlags = log.Ldate | log.Ltime
			initLogger(ioutil.Discard)
			Info.Println("deactivated debug mode")
		default:
			return errors.New("unknown mode:" + args[0])
		}
		return nil
	})
	RegisterCommand("help", "> help [command]", func(args []string) error {
		if len(args) == 0 {
			fmt.Println("available commands:")
			for _, v := range commands {
				fmt.Printf("%2s-%s\n", "", v.Name)
			}
		} else if len(args) == 1 {
			cmd, ok := commands[args[0]]
			if !ok {
				return errors.New("command not found")
			}
			fmt.Printf("usage of %s:\n%s\n", cmd.Name, cmd.Usage)
		} else {
			return errors.New("bad arguments")
		}
		return nil
	})
}

func runCLI() {
	for {
		fmt.Print(cmdSymbol)
		var input string
		fmt.Scanln(&input)

		tokens := strings.Split(input, " ")
		for i, _ := range tokens {
			tokens[i] = strings.Trim(tokens[i], " ")
		}

		if len(tokens) < 1 {
			continue
		} else {
			name := strings.ToLower(tokens[0])
			cmd, ok := commands[name]
			if !ok {
				fmt.Println("error: command not found")
				continue
			}

			err := cmd.Execute(tokens[1:])
			if err != nil {
				fmt.Println("error:", err)
			}
		}
	}
}
