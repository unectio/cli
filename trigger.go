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
	"strings"
)

var tgcol = apilet.FnTriggers

func doTrigger(cmd int, name *string) {
	tg_actions := map[int]func(*string){}

	tg_actions[CmdAdd] = triggerAdd
	tg_actions[CmdList] = triggerList
	tg_actions[CmdInfo] = triggerInfo
	tg_actions[CmdDel] = triggerDel

	doTargetCmd(cmd, name, tg_actions)
}

type elementTg struct{ *api.FuncTriggerImage }

func (te elementTg) id() string {
	return string(te.FuncTriggerImage.Id)
}

func (te elementTg) name() string {
	return te.FuncTriggerImage.Name
}

func (te elementTg) short() string {
	if te.FuncTriggerImage.URL != nil {
		return fmt.Sprintf("url: %s", te.FuncTriggerImage.URL.URL)
	}

	return ""
}

func (te elementTg) long() []*field {
	if te.FuncTriggerImage.URL != nil {
		return []*field{
			{
				name: "URL",
				data: te.FuncTriggerImage.URL.URL,
			},
			{
				name: "Auth",
				data: string(te.FuncTriggerImage.URL.AuthId),
			},
		}
	}
	if te.FuncTriggerImage.Cron != nil {
		return []*field{
			{
				name: "Tab",
				data: te.FuncTriggerImage.Cron.Tab,
			},
			{
				name: "Cron args",
				data: te.FuncTriggerImage.Cron.Args,
			},
		}
	}

	return nil
}

func triggerAdd(name *string) {
	var fn, tname string
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	var vfn = goopt.String([]string{"-f", "--function"}, "", "function name/id")
	var vsrc = goopt.String([]string{"-s", "--tsource"}, "", "trigger source")
	var vurl = goopt.String([]string{"-u", "--url"}, "", "trigger URL")
	var va = goopt.String([]string{"-a", "--auth"}, "", "URL trigger auth name/id")
	var vct = goopt.String([]string{"--crontab"}, "", "Cron trigger tab")
	var vca = goopt.String([]string{"--cronargs"}, "", "Cron trigger args in foo=bar:... format")
	goopt.Parse(nil)

	if strings.Contains(*name, "/") {
		x := strings.SplitN(*name, "/", 2)
		if len(x) != 2 {
			fatal("Specify function/trigger separated by \"/\" ")
		}
		fn = x[0]
		tname = x[1]
	} else {
		tname = *name
		fn = *vfn
	}

	fid := resolve(fcol, fn)

	tra := api.FuncTriggerImage{}
	tra.Name = generate(tname, "tg")

	switch *vsrc {
	case "url":
		tra.URL = &api.URLTrigImage{URL: api.AutoValue}
		if *va != "" {
			tra.URL.AuthId = resolve(authcol, *va)
		}
		if *vurl != "" {
			tra.URL.URL = api.URLProjectPfx + *vurl
		}
	case "cron":
		tra.Cron = &api.CronTrigImage{}
		if *vct != "" {
			tra.Cron.Tab = *vct
		}
		if *vca != "" {
			tra.Cron.Args = make(map[string]string)
			for _, a := range strings.Split(*vca, ":") {
				x := strings.SplitN(a, "=", 2)
				if len(x) != 2 {
					fatal("Bad cron arg %s", a)
				}
				tra.Cron.Args[x[0]] = x[1]
			}
		}
	}

	makeReq(tgcol.Sub(string(fid)).Add(&tra), &tra)

	showAddedElement(elementTg{&tra})
}

func triggerList(fcname *string) {
	var fn string
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2])
	goopt.ExtraUsage = ""
	var vfn = goopt.String([]string{"-f", "--function"}, "", "function name/id")
	goopt.Parse(nil)

	if fcname != nil && strings.HasPrefix(*fcname, "-") == false {
		fn = *fcname
	} else {
		fn = *vfn
	}

	fid := resolve(fcol, fn)

	var tgs []*api.FuncTriggerImage

	makeReq(tgcol.Sub(string(fid)).List(), &tgs)

	for _, tg := range tgs {
		showListElement(elementTg{tg})
	}
}

func triggerDel(name *string) {
	var fn, tname string
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	var vfn = goopt.String([]string{"-f", "--function"}, "", "function name/id")
	goopt.Parse(nil)

	if strings.Contains(*name, "/") {
		x := strings.SplitN(*name, "/", 2)
		if len(x) != 2 {
			fatal("Specify function/trigger separated by \"/\" ")
		}
		fn = x[0]
		tname = x[1]
	} else {
		tname = *name
		fn = *vfn
	}

	fnid := resolve(fcol, fn)
	tcol := tgcol.Sub(string(fnid))
	tgid := resolve(tcol, tname)

	makeReq(tcol.Delete(string(tgid)), nil)
}

func triggerInfo(name *string) {
	var fn, tname string
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	var vfn = goopt.String([]string{"-f", "--function"}, "", "function name/id")
	goopt.Parse(nil)

	if strings.Contains(*name, "/") {
		x := strings.SplitN(*name, "/", 2)
		if len(x) != 2 {
			fatal("Specify function/trigger separated by \"/\" ")
		}
		fn = x[0]
		tname = x[1]
	} else {
		tname = *name
		fn = *vfn
	}

	fnid := resolve(fcol, fn)
	tcol := tgcol.Sub(string(fnid))
	tgid := resolve(tcol, tname)

	var tg api.FuncTriggerImage

	makeReq(tcol.Info(string(tgid)), &tg)

	showInfoElement(elementTg{&tg})
}
