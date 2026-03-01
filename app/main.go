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
	"cd":   "builtin",
}

// this is a where you parse the args
func parseArgs(line string) []string {
	var args []string
	var b strings.Builder
	inSingle := false

	flush := func() {
		if b.Len() > 0 {
			args = append(args, b.String())
			b.Reset()
		}
	}

	// go through all the characters in the input
	for _, r := range line {
		switch {
		case r == '\'':
			// this toggles in and out of quotes to see if we are inside or not
			inSingle = !inSingle
		// i check here if i hit whitespace and we are not in single quotes
		case (r == ' ' || r == '\t') && !inSingle:
			flush()

		default:
			// add the normal character to the token
			b.WriteRune(r)
		}
	}
	flush()
	return args
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
			// here i have to deal with the quotes of the strings
			args := cmd[len("echo "):]
			// find where there are single quotes and remove it
			if ok := strings.Contains(args, "'"); ok {
				args = strings.ReplaceAll(args, "'", "")
				fmt.Println(args)

			} else {
				curr := strings.Fields(args)
				fmt.Println(strings.Join(curr, " "))
			}
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
		} else if strings.HasPrefix(cmd, "cd ") {
			// basically i check what comes after it and i check if the dir exists
			// if it exists i go into the dir
			args := strings.TrimSpace(cmd[len("cd "):])
			if args == "~" || args == "" {
				// get the home directory
				home, err := os.UserHomeDir()
				if err != nil {
					fmt.Fprintf(os.Stderr, "cd: %v\n", err)
				} else {
					if err := os.Chdir(home); err != nil {
						fmt.Fprintf(os.Stderr, "cd: %s: %v\n", args, err)
					}
				}
				continue
			} else {
				// checking the dir
				if err := os.Chdir(args); err != nil {
					if os.IsNotExist(err) {
						fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", args)
					} else {
						fmt.Fprintf(os.Stderr, "cd: %s: %v\n", args, err)
					}
					continue
				}
			}
			continue
		}
		res := parseArgs(cmd)
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
