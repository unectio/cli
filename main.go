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
)

const (
	defaultConfig string = "/etc/lets.config"
)

func fatal(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func usage_summary() {
	fmt.Printf("Usage: %s <command> [<target>] [<object>] [<options>]\n", os.Args[0])

	usage_commands()
	usage_targets()

	fmt.Printf("Try '%s <command> --help' for command help\n", os.Args[0])
	fmt.Printf("    '%s <command> <target> --help' for target help\n", os.Args[0])
	fmt.Printf("    '%s <command> <target> <object> --help' for help with object arguments\n", os.Args[0])
	fmt.Printf("\nDefault configuration file is %s\n\n", defaultConfig)
}

func usage_commands_string() string {
	usage := fmt.Sprintf("Commands (<tgt> means command requires a target):\n")
	cmds := listCommands()
	for _, cmd := range cmds {
		usage += fmt.Sprintf("\t%s\n", cmd)
	}
	return usage
}

func usage_commands() {
	fmt.Printf("Available commands (<tgt> means command requires a target):\n")
	cmds := listCommands()
	for _, cmd := range cmds {
		fmt.Printf("\t%s\n", cmd)
	}
}

func usage_targets() {
	fmt.Printf("Recognized targets:\n")
	tgts := listTargets()
	for _, tgt := range tgts {
		fmt.Printf("\t%s\n", tgt)
	}
}

func usage_targets_string() string {
	usage := fmt.Sprintf("Recognized targets:\n")
	tgts := listTargets()
	for _, tgt := range tgts {
		usage += fmt.Sprintf("\t%s\n", tgt)
	}
	return usage
}

func usage_footer_string() string {
	usage := fmt.Sprintf("\nTry '%s <command> --help' for command help\n", os.Args[0])
	usage += fmt.Sprintf("    '%s <command> <target> --help' for target help\n", os.Args[0])
	usage += fmt.Sprintf("    '%s <command> <target> <object> --help' for help with object arguments\n", os.Args[0])
	usage += fmt.Sprintf("Default configuration file is %s", defaultConfig)
	return usage
}

func usage_short_string() string {
	return fmt.Sprintf("Usage: %s <command> [<target>] [<object>] [<options>]\n", os.Args[0])
}

func usage_string() string {
	usage := ""
	if goopt.Summary != "" {
		usage += fmt.Sprintf("%s", goopt.Summary)
	}
	usage += fmt.Sprintf("\n%s", goopt.Help())
	if goopt.ExtraUsage != "" {
		usage += fmt.Sprintf("%s\n", goopt.ExtraUsage)
	}
	return usage
}

var debug = goopt.Flag([]string{"-d", "--debug"}, []string{}, "Pring debugging info", "")
var dryrun = goopt.Flag([]string{"--dry-run"}, []string{}, "Do not do requests for real", "")

func main() {
	goopt.Usage = usage_string
	goopt.Summary = usage_short_string() + usage_commands_string() + "\n" + usage_targets_string()
	goopt.ExtraUsage = usage_footer_string()
	goopt.Version = "1.4"

	if len(os.Args) <= 1 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		fmt.Println(goopt.Usage())
		os.Exit(1)
	}

	/* Usage is always $ lets <command> [<object>] [<options>] */
	c := os.Args[1]
	//	os.Args = os.Args[1:]

	cmd := getCommand(c)
	cmd.Do()
}
