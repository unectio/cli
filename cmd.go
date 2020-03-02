/////////////////////////////////////////////////////////////////////////////////
//
// Copyright (C) 2019-2020, Unectio Inc, All Right Reserved.
//
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
/////////////////////////////////////////////////////////////////////////////////

package main

import (
	"os"
	"fmt"
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
	case "add", "create":
		return &Command { Do: doAdd }
	case "del", "delete":
		return &Command { Do: doDel }
	case "ls", "list":
		return &Command { Do: doList }
	case "upd", "update", "set":
		return &Command { Do: doUpdate }
	case "info", "get":
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
		"add | create *",
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

	if len(os.Args) <= 1 || os.Args[1] == "--help" {
		fmt.Printf("Specify a target\n")
		usage_targets()
		return
	}

	cmd := os.Args[0]
	t := os.Args[1]
	os.Args = os.Args[1:]

	if c != CmdList {
		if len(os.Args) <= 1 {
			fmt.Print("Specify an object name or %id for operations with existing objects\n")
			return
		}
		name = &os.Args[1]
		os.Args = os.Args[1:]
	}

	tgt := getTarget(t)
	tgt.Do(c, name)
}
