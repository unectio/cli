package main

import (
	"os"
	"fmt"
	"flag"
)

const (
	CmdAdd = 1
	CmdDel = 2
	CmdList = 3
	CmdUpdate = 4
	CmdInfo = 5
)

type Command struct {
	Do	func()
}

func getCommand(c string) *Command {
	switch c {
	case "add":
		return &Command { Do: doAdd }
	case "del", "delete":
		return &Command { Do: doDel }
	case "ls", "list":
		return &Command { Do: doList }
	case "upd", "update":
		return &Command { Do: doUpdate }
	case "info":
		return &Command { Do: doInfo }
	case "run":
		return &Command { Do: functionRun }
	case "load":
		return &Command { Do: loadSpec }
	}

	return &Command { Do: func() {
		fmt.Printf("Unknown command %s\n", c)
		usage_commands()
	} }
}

func listCommands() []string {
	return []string {
		"ls  | list   *",
		"add          *",
		"del | delete *",
		"upd | update *",
		"info         *",
		"run          (function)",
		"load         (spec file)",
	}
}

func doAdd() { doTarget(CmdAdd) }
func doDel() { doTarget(CmdDel) }
func doList() { doTarget(CmdList) }
func doUpdate() { doTarget(CmdUpdate) }
func doInfo() { doTarget(CmdInfo) }

func doTarget(c int) {
	var name *string

	if len(os.Args) <= 1 || os.Args[1] == "-help" {
		fmt.Printf("Specify a target\n")
		usage_targets()
		return
	}

	t := os.Args[1]
	os.Args = os.Args[1:]

	if c != CmdList {
		name = flag.String("n", "", "name or %id")
	}

	tgt := getTarget(t)
	tgt.Do(c, name)
}
