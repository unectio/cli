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
)

type element interface {
	id() string
	name() string
	short() string
	long() []*field
}

type field struct {
	name string
	data interface{}
}

func showListElement(le element) {
	fmt.Printf("%s: %-32s (%s)\n", le.id(), le.name(), le.short())
}

func showInfoElement(le element) {
	fmt.Printf("%-12s: %s\n", "Id", le.id())
	fmt.Printf("%-12s: %s\n", "Name", le.name())
	for _, f := range le.long() {
		fmt.Printf("HHH%-12s: %v\n", f.name, f.data)
	}
}

func showAddedElement(le element) {
	fmt.Printf("Added, id=%s (%s)\n", le.id(), le.short())
}

type fmtMap map[string]string

func (m fmtMap) String() string {
	ret := ""
	for k, v := range m {
		ret += k + "=" + v + ":"
	}
	return ret
}

type fmtArray []string

func (m fmtArray) String() string {
	ret := ""
	for _, x := range m {
		ret += ", " + x
	}
	if len(ret) > 2 {
		return ret[2:]
	}

	return "<empty>"
}
