package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// create the slice of strings
var commands = map[string]string{
	"exit": "builtin",
	"echo": "builtin",
	"type": "builtin",
	"pwd":  "builtin",
}

func main() {
	// read the input
	reader := bufio.NewReader(os.Stdin)
	for {
		// so now for echo i have to print the words after it
		fmt.Print("$ ")
		cmd, err := reader.ReadString('\n')
		//   if its an error with end of line and since its still standard invalid cases[all for now]
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				return
			}
			fmt.Fprintf(os.Stderr, "read error %v\n", err)
			return
		}
		// read the input
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		if cmd == "pwd" {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "pwd: %v\n", err)
			} else {
				fmt.Println(dir)
			}
			continue
		}
		if cmd == "exit" {
			break
		} else if strings.HasPrefix(cmd, "echo ") {
			fmt.Println(cmd[len("echo "):])
			continue
		} else if strings.HasPrefix(cmd, "type ") {
			// get the first 4 letters after it
			name := strings.TrimSpace(cmd[len("type "):])
			if val, ok := commands[name]; ok {
				fmt.Println(name + " is a shell " + val)
			} else {
				// this is the point where i check my path since its not a built in command
				path, err := exec.LookPath(name)
				if err != nil {
					fmt.Println(name + ": not found")
				} else {
					fmt.Println(name + " is " + path)
				}
			}
			continue
		}
		res := strings.Fields(cmd)
		if len(res) == 0 {
			continue
		}
		// get program and its args
		program := res[0]
		args := res[1:]
		path, err := exec.LookPath(program)
		if err != nil {
			fmt.Printf("%s: command not found\n", program)
			continue
		}

		// run the executable and add the args
		// i had to ensure that path wasnt added to the args so i specifically passed the prog inside
		c := exec.Command(program, args...)
		c.Path = path
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		_ = c.Run()

	}
}
