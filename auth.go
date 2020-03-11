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
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
	"os"
)

var authcol = apilet.AuthMethods

func doAuth(cmd int, name *string) {
	a_actions := map[int]func(namep *string){}
	a_actions[CmdAdd] = authAdd
	a_actions[CmdList] = authList
	a_actions[CmdInfo] = authInfo
	a_actions[CmdDel] = authDel

	doTargetCmd(cmd, name, a_actions)
}

func authAdd(name *string) {
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	var key = goopt.String([]string{"-k", "--key"}, "", "jwt key (base64-encoded or auto)")
	goopt.Parse(nil)

	ai := api.AuthMethodImage{}
	ai.Name = generate(*name, "am")

	if *key != "" {
		ai.JWT = &api.AuthJWTImage{Key: *key}
	}

	makeReq(authcol.Add(&ai), &ai)

	fmt.Printf("Added auth method (id %s)\n", ai.Id)
}

func authDel(name *string) {
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	goopt.Parse(nil)

	aid := resolve(authcol, *name)

	makeReq(authcol.Delete(string(aid)), nil)
}

func authList(_ *string) {
	var as []*api.AuthMethodImage
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2])
	goopt.ExtraUsage = ""
	goopt.Parse(nil)

	makeReq(authcol.List(), &as)

	for _, am := range as {
		fmt.Printf("%s: %s\n", am.Id, am.Name)
	}
}

func authInfo(name *string) {
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	goopt.Parse(nil)

	aid := resolve(authcol, *name)

	var ai api.AuthMethodImage

	makeReq(authcol.Info(string(aid)), &ai)

	fmt.Printf("Id:             %s\n", ai.Id)
	fmt.Printf("Name:           %s\n", ai.Name)
	if ai.JWT != nil {
		fmt.Printf("Key:            %s\n", ai.JWT.Key)
	}
}
