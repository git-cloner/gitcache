package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func execShell(cmd string, args []string) string {
	var argss = ""
	for i := 0; i < len(args); i++ {
		argss = argss + args[i] + " "
	}
	log.Println(cmd + " " + string(argss))
	var command = exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	var err = command.Start()
	if err != nil {
		return err.Error()
	}
	err = command.Wait()
	if err != nil {
		return err.Error()
	}
	return ""
}

func main() {
	log.Printf("gitcache client\n")
	args := os.Args
	var isClone = false
	for i := 0; i < len(args); i++ {
		if strings.Contains(args[i], "clone") {
			isClone = true
			break
		}
	}
	var isDepth = false
	for i := 0; i < len(args); i++ {
		if isClone && strings.Contains(args[i], "https://github.com") {
			args[i] = strings.Replace(args[i], "https://github.com", "https://gitclone.com/github.com", -1)
		}
		if strings.Contains(args[i], "depth") {
			isDepth = true
		}
	}
	if isClone && (!isDepth) {
		args = append(args, "--depth=1")
	}
	execShell("git", args[1:])
}
