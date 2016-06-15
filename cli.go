package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Command is a CLI command.
type Command struct {
	Name, Usage string
	Execute     func([]string) error
}

const (
	cmdSymbol = "~ "
)

var (
	commands []*Command
)

// GetCommandNames returns a slice of all command names.
func GetCommandNames() []string {
	allNames := make([]string, 0, len(commands))
	for _, cmd := range commands {
		allNames = append(allNames, cmd.Name)
	}
	return allNames
}

// GetCommand either returns a command or and error.
func GetCommand(name string) (*Command, error) {
	for _, cmd := range commands {
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
	commands = append(commands, &cmd)
}

// init loads a standard set of command.
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
	RegisterCommand("build", "> build", func(args []string) error {
		const routeBaseDir = "build"
		var paths = []string{"/"}
		for _, p := range Posts {
			paths = append(paths, p.GetURL())
		}
		for _, p := range Pages {
			paths = append(paths, p.GetURL())
		}
		blogBasePath := "http://localhost:" + strconv.Itoa(Config.Server.Port)
		err := os.Mkdir(routeBaseDir, 0700)
		if err != nil {
			return err
		}
		for _, p := range paths {
			url := blogBasePath + p
			fmt.Println("Loading " + url)
			resp, err := http.Get(url)
			if err != nil {
				return err
			}

			defer resp.Body.Close()
			if p == "/" {
				p = "/index"
				url = blogBasePath + "/index"
			}
			baseDir := filepath.Dir(filepath.Join(routeBaseDir, p))
			err = os.MkdirAll(baseDir, 0700)
			if err != nil {
				return err
			}

			file, err := os.Create(filepath.Join(baseDir, filepath.Base(url)+".html"))
			if err != nil {
				return err
			}

			defer file.Close()
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func check(err error) {
	if err != nil {
		panic(err)
	}
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
		for i := range tokens {
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
