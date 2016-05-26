package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func runCLI() {
	for {
		fmt.Print("> ")
		var command, arg string
		fmt.Scanln(&command, &arg)

		switch command {
		case "reload":
			// Reload all templates and posts
			LoadTemplates()
			LoadPosts()
			LoadPages()
			Info.Println("Reloaded posts and templates")
		case "stop":
			// Stops the server
			Info.Println("Forcing shutdown")
			os.Exit(1)
		case "debug":
			switch arg {
			case "on":
				logFlags = log.Ltime | log.Lshortfile
				initLogger(os.Stdout)
				Info.Println("Activated debug mode")
			case "off":
				logFlags = log.Ldate | log.Ltime
				initLogger(ioutil.Discard)
				Info.Println("Deactivated debug mode")
			}
		case "help":
			fmt.Println("Print help here.")
		}
	}
}
