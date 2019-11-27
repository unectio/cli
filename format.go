package main

import (
	"fmt"
)

type element interface {
	id()	string
	name()	string
	short()	string
	long()	[]*field
}

type field struct {
	name	string
	data	interface{}
}

func showListElement(le element) {
	fmt.Printf("%s: %-32s (%s)\n", le.id(), le.name(), le.short())
}

func showInfoElement(le element) {
	fmt.Printf("%-12s: %s\n", "Name", le.name())
	for _, f := range le.long() {
		fmt.Printf("%-12s: %v\n", f.name, f.data)
	}
}

func showAddedElement(le element) {
	fmt.Printf("Added, id=%s (%s)\n", le.id(), le.short())
}

type fmtMap map[string]string

func (m fmtMap)String() string {
	ret := ""
	for k, v := range m {
		ret += k + "=" + v + ":"
	}
	return ret
}

type fmtArray []string

func (m fmtArray)String() string {
	ret := ""
	for _, x := range m {
		ret += ", " + x
	}
	if len(ret) > 2 {
		return ret[2:]
	}

	return "<empty>"
}
