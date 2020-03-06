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
	"encoding/json"
	"flag"
	"fmt"
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
	rq "github.com/unectio/util/request"
	"os"
	"strings"
)

var fcol = apilet.Functions

func doFunction(cmd int, name *string) {
	fn_actions := map[int]func(name *string){}

	fn_actions[CmdAdd] = functionAdd
	fn_actions[CmdDel] = functionDelete
	fn_actions[CmdList] = functionList
	fn_actions[CmdUpdate] = functionUpdate
	fn_actions[CmdInfo] = functionInfo

	doTargetCmd(cmd, name, fn_actions)
}

type elementFn struct{ *api.FunctionImage }

func (fe elementFn) id() string {
	return string(fe.FunctionImage.Id)
}

func (fe elementFn) name() string {
	return fe.FunctionImage.Name
}

func (fe elementFn) short() string {
	return ""
}

func (fe elementFn) long() []*field {
	return []*field{
		{
			name: "State",
			data: fe.FunctionImage.State,
		},
		{
			name: "Env",
			data: fmtArray(fe.FunctionImage.Env),
		},
		{
			name: "Code balancer",
			data: fe.FunctionImage.CodeBalancer,
		},
	}
}

func functionAdd(name *string) {
	var env string
	const (
		default_value = ""
		usage         = "environment (key=val;...)"
	)
	flag.StringVar(&env, "environment", default_value, usage)
	flag.StringVar(&env, "e", default_value, usage+" (shorthand)")
	flag.Parse()

	fa := api.FunctionImage{}
	fa.Name = generate(*name, "fn")

	if env != "" {
		fa.Env = parseEnv(env)
	}

	makeReq(fcol.Add(&fa), &fa)

	showAddedElement(elementFn{&fa})
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
	var inf string
	const (
		default_value = ""
		usage         = "what to show (logs, stats)"
	)
	flag.StringVar(&inf, "information", default_value, usage)
	flag.StringVar(&inf, "i", default_value, usage+" (shorthand)")
	flag.Parse()

	switch inf {
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

	var cis []*api.CodeImage

	makeReq(ccol.Sub(string(fnid)).List(), &cis)

	for _, ci := range cis {
		makeReq(ccol.Sub(string(fnid)).Delete(elementCode{ci}.id()), nil)
	}

	makeReq(fcol.Delete(string(fnid)), nil)
}

func functionRun() {

	if len(os.Args) <= 1 {
		fatal("Specify function/code to run")
	}

	var name, code, req string

	if strings.Contains(os.Args[1], "/") {
		x := strings.SplitN(os.Args[1], "/", 2)
		if len(x) != 2 {
			fatal("Specify function/code to run separated by \"/\" ")
		}
		name = x[0]
		code = x[1]
	} else {
		code = os.Args[1]
		const (
			fndefault_value = ""
			fnusage         = "function name/id"
		)
		flag.StringVar(&name, "function", fndefault_value, fnusage)
		flag.StringVar(&name, "f", fndefault_value, fnusage+" (shorthand)")
	}
	os.Args = os.Args[1:]

	const (
		rdefault_value = ""
		rusage         = "request (JSON string)"
	)
	flag.StringVar(&req, "rq", rdefault_value, rusage)
	flag.StringVar(&req, "request", rdefault_value, rusage+" (shorthand)")
	flag.Parse()

	var rreq api.FuncRun
	var res api.RunResponse

	err := json.Unmarshal([]byte(req), &rreq.Req)
	if err != nil {
		fatal("Bad req param: " + err.Error())
	}

	fnid := resolve(fcol, name)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, code)

	makeReq(rq.Req("", "functions/"+string(fnid)+"/code/"+string(cver)+"/run").B(&rreq), &res)

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
	var env string
	const (
		default_value = ""
		usage         = "environment (key=val;...)"
	)
	flag.StringVar(&env, "environment", default_value, usage)
	flag.StringVar(&env, "e", default_value, usage+" (shorthand)")
	flag.Parse()

	fnid := resolve(fcol, *name)

	switch {
	case env != "":
		functionUpdateEnv(fnid, env)
	}
}

var envprop = apilet.FnEnvironment

func functionUpdateEnv(fnid api.ObjectId, env string) {
	envi := parseEnv(env)

	makeReq(envprop.Set(string(fnid), envi), nil)
}
