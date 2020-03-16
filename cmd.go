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
	"fmt"
	goopt "github.com/droundy/goopt"
	"os"
	"strings"
)

const (
	CmdAdd    = 1
	CmdDel    = 2
	CmdList   = 3
	CmdUpdate = 4
	CmdInfo   = 5
)

type Command struct {
	Do func()
}

func getCommand(c string) *Command {
	switch c {
	case "add", "create":
		return &Command{Do: doAdd}
	case "del", "delete":
		return &Command{Do: doDel}
	case "ls", "list":
		return &Command{Do: doList}
	case "upd", "update", "set":
		return &Command{Do: doUpdate}
	case "info", "get":
		return &Command{Do: doInfo}
	case "run":
		return &Command{Do: functionRun}
	case "load":
		return &Command{Do: loadSpec}
	case "pull":
		return &Command{Do: repoPull}
	}

	return &Command{Do: func() {
		goopt.Summary = fmt.Sprintf("Unknown command \"%s\"\n\n", c) + usage_commands_string()
		goopt.ExtraUsage = ""
		fmt.Println(goopt.Usage())
	}}
}

func listCommands() []string {
	return []string{
		"ls   | list   <tgt>",
		"add  | create <tgt>",
		"del  | delete <tgt>",
		"upd  | update <tgt>",
		"info | get    <tgt>",
		"run          (function)",
		"pull         (repository)",
		"load         (spec file)",
	}
}

func doAdd()    { doTarget(CmdAdd) }
func doDel()    { doTarget(CmdDel) }
func doList()   { doTarget(CmdList) }
func doUpdate() { doTarget(CmdUpdate) }
func doInfo()   { doTarget(CmdInfo) }

func doTarget(c int) {
	var name *string

	if len(os.Args) <= 2 || os.Args[2] == "--help" || os.Args[2] == "-h" {
		goopt.Summary = fmt.Sprintf("Specify a target\n\n") + usage_targets_string()
		goopt.ExtraUsage = ""
		fmt.Println(goopt.Usage())
		return
	}

	t := os.Args[2]

	if c != CmdList {
		if len(os.Args) <= 3 || strings.HasPrefix(os.Args[3], "-") {
			goopt.Summary = fmt.Sprint("Specify an object name or %id \n")
			goopt.ExtraUsage = ""
			fmt.Println(goopt.Usage())
			return
		}
		name = &os.Args[3]
	} else {
		if len(os.Args) > 3 {
			name = &os.Args[3]
		}
	}

	tgt := getTarget(t)
	tgt.Do(c, name)
}
