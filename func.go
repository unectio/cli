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
	"flag"
	"strings"
	"github.com/unectio/api"
	"io/ioutil"
	"encoding/json"
	"github.com/unectio/api/apilet"
	"github.com/unectio/util/request"
)

var fcol = apilet.Functions

func doFunction(cmd int, name *string) {
	fn_actions := map[int]func(name *string) {}

	fn_actions[CmdAdd] = functionAdd
	fn_actions[CmdDel] = functionDelete
	fn_actions[CmdList] = functionList
	fn_actions[CmdUpdate] = functionUpdate
	fn_actions[CmdInfo] = functionInfo

	doTargetCmd(cmd, name, fn_actions)
}

type elementFn struct { *api.FunctionImage }

func (fe elementFn)id() string {
	return string(fe.FunctionImage.Id)
}

func (fe elementFn)name() string {
	return fe.FunctionImage.Name
}

func (fe elementFn)short() string {
	return ""
}

func (fe elementFn)long() []*field {
	return []*field {
		{
			name:	"State",
			data:	fe.FunctionImage.State,
		},
		{
			name:	"Env",
			data:	fmtArray(fe.FunctionImage.Env),
		},
		{
			name:	"Code balancer",
			data:	fe.FunctionImage.CodeBalancer,
		},
	}
}

func functionAdd(name *string) {
	env := flag.String("e", "", "environment (key=val;...)")
	flag.Parse()

	fa := api.FunctionImage{}
	fa.Name = generate(*name, "fn")

	if *env != "" {
		fa.Env = parseEnv(*env)
	}

	makeReq(fcol.Add(&fa), &fa)

	showAddedElement(elementFn{&fa})
}

func parseCode(src string, ci *api.SourceImage) {
	if strings.HasPrefix(src, "http") {
		ci.URL = src
	} else if strings.HasPrefix(src, "repo:") {
		x := strings.SplitN(src, ":", 3)
		ci.RepoId = api.ObjectId(x[1])
		ci.Path = x[2]
	} else {
		var err error

		ci.Text, err = ioutil.ReadFile(src)
		if err != nil {
			fatal("Error reading sources: %s\n", err.Error())
		}
	}
}

func parseEnv(envs string) []string {
	return strings.Split(envs, ";")
}

func functionList(_ *string) {
	var fns []*api.FunctionImage

	flag.Parse()
	makeReq(fcol.List(), &fns)

	for _, fn := range fns {
		showListElement(elementFn{fn})
	}
}

func functionInfo(name *string) {
	inf := flag.String("i", "", "what to show (logs, stats)")
	flag.Parse()

	switch *inf {
	case "stats":
		functionStats(name)
	case "logs":
		functionLogs(name)
	default:
		functionCommonInfo(name)
	}
}

func functionCommonInfo(name *string) {
	var fn api.FunctionImage

	fnid := resolve(fcol, *name)
	makeReq(fcol.Info(string(fnid)), &fn)

	showInfoElement(elementFn{&fn})
}

func functionDelete(name *string) {
	flag.Parse()
	fnid := resolve(fcol, *name)
	makeReq(fcol.Delete(string(fnid)), nil)
}

func functionRun() {
	name := flag.String("n", "", "function name")
	code := flag.String("c", "", "code name")
	req := flag.String("rq", "", "request (JSON string)")
	flag.Parse()

	var rreq api.FuncRun
	var res api.RunResponse

	err := json.Unmarshal([]byte(*req), &rreq.Req)
	if err != nil {
		fatal("Bad req param: " + err.Error())
	}

	fnid := resolve(fcol, *name)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, *code)

	makeReq(rq.Req("", "functions/" + string(fnid) + "/code/" + string(cver) + "/run").B(&rreq), &res)

	fmt.Printf("Status:            %d\n", res.Status)
	fmt.Printf("Time taken:        %dus\n", res.LatUs)
	fmt.Printf("Returned:          %s\n", string(res.Res))

	if res.Out != "" {
		fmt.Printf("Out:               %s\n", res.Out)
	}
	if res.Err != "" {
		fmt.Printf("Err:               %s\n", res.Err)
	}
}

func functionUpdate(name *string) {
	env := flag.String("e", "", "environment (key=val:...)")
	flag.Parse()

	fnid := resolve(fcol, *name)

	switch {
	case *env != "":
		functionUpdateEnv(fnid, *env)
	}
}

var envprop = apilet.FnEnvironment

func functionUpdateEnv(fnid api.ObjectId, env string) {
	envi := parseEnv(env)

	makeReq(envprop.Set(string(fnid), envi), nil)
}
