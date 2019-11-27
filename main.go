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

func usage() {
	fmt.Printf("Usage: %s <command> [<target>] [<options>]\n", os.Args[0])
	fmt.Printf("Commands (* means cmd is targeted):\n")
	cmds := listCommands()
	for _, cmd := range cmds {
		fmt.Printf("\t%s\n", cmd)
	}

	fmt.Printf("Targets:\n")
	tgts := listTargets()
	for _, tgt := range tgts {
		fmt.Printf("\t%s\n", tgt)
	}

	os.Exit(1)
}

var debug *bool
var dryrun *bool

func main() {
	if len(os.Args) <= 2 {
		usage()
	}

	debug = flag.Bool("debug", false, "Print debugging info")
	dryrun = flag.Bool("dry-run", false, "Do not do requests for real")

	/* Usage is always $ lets <command> [<object>] [<options>] */
	c := os.Args[1]
	os.Args = os.Args[1:]

	cmd := getCommand(c)
	cmd.Do()
}
