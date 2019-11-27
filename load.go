// +build batcher

package main

import (
	"os"
	"flag"
	"this/batcher"
)

func loadSpec() {
	name := flag.String("f", "", "spec file name")
	flag.Parse()

	l, err := getLogin()
	if err != nil {
		fatal("Cannot login: %s", err.Error())
	}

	/*
	 * XXX -- if possible -- use the platform call
	 */

	f, err := os.Open(*name)
	if err != nil {
		fatal("Cannot open file: %s", err.Error())
	}

	defer f.Close()

	err = batcher.Process(f, l)
	if err != nil {
		fatal("Cannot process file: %s", err.Error())
	}
}
