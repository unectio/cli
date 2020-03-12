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
)

type Target struct {
	Do func(cmd int, namep *string)
}

func getTarget(t string) *Target {
	switch t {
	case "fn", "func", "function":
		return &Target{Do: doFunction}
	case "code":
		return &Target{Do: doCode}
	case "rt", "router":
		return &Target{Do: doRouter}
	case "repo", "repository":
		return &Target{Do: doRepo}
	case "sec", "secret":
		return &Target{Do: doSecret}
	case "tg", "trig", "trigger":
		return &Target{Do: doTrigger}
	case "am", "auth_method":
		return &Target{Do: doAuth}
	}

	return &Target{Do: func(_ int, _ *string) {
		goopt.Summary = fmt.Sprintf("Unknown target \"%s\"\n\n", t) + usage_targets_string()
		goopt.ExtraUsage = ""
		fmt.Println(goopt.Usage())
	}}
}

func listTargets() []string {
	return []string{
		"fn | func | function",
		"code",
		"rt | router",
		"repo | repository",
		"sec | secret",
		"tg | trig | trigger",
		"am | auth_method",
	}
}

func doTargetCmd(cmd int, namep *string, actions map[int]func(namep *string)) {
	fn, ok := actions[cmd]
	if !ok {
		fn = func(_ *string) {
			goopt.Summary = usage_targets_string()
			goopt.ExtraUsage = ""
			fmt.Println(goopt.Usage())
		}
	}
	fn(namep)
}
