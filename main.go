package main

import (
	"os"
	"fmt"
	"flag"
)

const (
	defaultConfig string = "/etc/lets.config"
)

func fatal(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg + "\n", args...)
	os.Exit(1)
}

func usage_summary() {
	fmt.Printf("Usage: %s <command> [<target>] [<options>]\n", os.Args[0])

	usage_commands()
	usage_targets()

	fmt.Printf("Try '%s <command>' for command help\n", os.Args[0])
	fmt.Printf("    '%s <command> <target>' for target help\n", os.Args[0])
	fmt.Printf("\nDefault configuration file is %s\n\n", defaultConfig)
}

func usage_commands() {
	fmt.Printf("Commands (* means cmd is targeted):\n")
	cmds := listCommands()
	for _, cmd := range cmds {
		fmt.Printf("\t%s\n", cmd)
	}
}

func usage_targets() {
	fmt.Printf("Targets:\n")
	tgts := listTargets()
	for _, tgt := range tgts {
		fmt.Printf("\t%s\n", tgt)
	}
}

var debug *bool
var dryrun *bool

func main() {
	if len(os.Args) <= 1 || os.Args[1] == "-help" {
		usage_summary()
		os.Exit(1)
	}

	debug = flag.Bool("debug", false, "Print debugging info")
	dryrun = flag.Bool("dry-run", false, "Do not do requests for real")

	/* Usage is always $ lets <command> [<object>] [<options>] */
	c := os.Args[1]
	os.Args = os.Args[1:]

	cmd := getCommand(c)
	cmd.Do()
}
