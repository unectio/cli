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
	"io/ioutil"
	"os"
	"strings"
)

var ccol = apilet.FnCodes

func doCode(cmd int, name *string) {
	co_actions := map[int]func(*string){}

	co_actions[CmdAdd] = codeAdd
	co_actions[CmdList] = codeList
	co_actions[CmdDel] = codeDel
	co_actions[CmdUpdate] = codeUpdate
	co_actions[CmdInfo] = codeInfo

	doTargetCmd(cmd, name, co_actions)
}

type elementCode struct{ *api.CodeImage }

func (ce elementCode) id() string {
	return string(ce.CodeImage.Id)
}

func (ce elementCode) name() string {
	return ce.CodeImage.Name
}

func (ce elementCode) short() string {
	return ce.CodeImage.Lang
}

func (ce elementCode) long() []*field {
	return []*field{
		{
			name: "State",
			data: ce.CodeImage.State,
		},
		{
			name: "Gen",
			data: ce.CodeImage.Gen,
		},
		{
			name: "Lang",
			data: ce.CodeImage.Lang,
		},
		{
			name: "Weight",
			data: ce.CodeImage.Weight,
		},
	}
}

func codeAdd(fcname *string) {
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	var fn, cname, lang, src string
	var w int
	var vfn = goopt.String([]string{"-f", "--function"}, "", "function name/id")
	var vlang = goopt.String([]string{"-l", "--language"}, "", "code language")
	var vsrc = goopt.String([]string{"-s", "--source"}, "", "sources (file name or url or repo:<repo name>:path)")
	var vw = goopt.Int([]string{"-w", "--weight"}, 0, "code weight")
	goopt.Parse(nil)

	if strings.Contains(*fcname, "/") {
		x := strings.SplitN(*fcname, "/", 2)
		if len(x) != 2 {
			fatal("Specify function/code separated by \"/\" ")
		}
		fn = x[0]
		cname = x[1]
	} else {
		cname = *fcname
		fn = *vfn
	}
	lang = *vlang
	src = *vsrc
	w = *vw

	fid := resolve(fcol, fn)

	var ci api.CodeImage

	ci.Name = generate(cname, "code")
	ci.Lang = lang
	ci.Weight = w
	ci.Source = &api.SourceImage{}
	parseCode(src, ci.Source)

	makeReq(ccol.Sub(string(fid)).Add(&ci), &ci)

	showAddedElement(elementCode{&ci})
}

func codeList(fcname *string) {
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

	var cis []*api.CodeImage

	makeReq(ccol.Sub(string(fid)).List(), &cis)

	for _, ci := range cis {
		showListElement(elementCode{ci})
	}
}

func codeDel(ver *string) {
	var fn, cname string
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2])
	goopt.ExtraUsage = ""
	var vfn = goopt.String([]string{"-f", "--function"}, "", "function name/id")
	goopt.Parse(nil)

	if strings.Contains(*ver, "/") {
		x := strings.SplitN(*ver, "/", 2)
		if len(x) != 2 {
			fatal("Specify function/code separated by \"/\" ")
		}
		fn = x[0]
		cname = x[1]
	} else {
		cname = *ver
		fn = *vfn
	}

	fnid := resolve(fcol, fn)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, cname)

	makeReq(xcol.Delete(string(cver)), nil)
}

func codeUpdate(ver *string) {
	var fn, cname, src string
	var w int
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	var vfn = goopt.String([]string{"-f", "--function"}, "", "function name/id")
	var vsrc = goopt.String([]string{"-s", "--source"}, "", "sources (file name or url or repo:<repo name>:path)")
	var vw = goopt.Int([]string{"-w", "--weight"}, 0, "code weight")
	goopt.Parse(nil)

	if strings.Contains(*ver, "/") {
		x := strings.SplitN(*ver, "/", 2)
		if len(x) != 2 {
			fatal("Specify function/code separated by \"/\" ")
		}
		fn = x[0]
		cname = x[1]
	} else {
		cname = *ver
		fn = *vfn
	}

	src = *vsrc
	w = *vw

	var ci api.CodeImage

	ci.Weight = w
	ci.Source = &api.SourceImage{}
	parseCode(src, ci.Source)

	fnid := resolve(fcol, fn)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, cname)

	makeReq(xcol.Upd(string(cver), &ci), nil)
}

func codeInfo(ver *string) {
	var fn, cname string
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	var vfn = goopt.String([]string{"-f", "--function"}, "", "function name/id")
	var voc = goopt.Flag([]string{"-C", "--code"}, []string{}, "Show code only", "")
	goopt.Parse(nil)

	if strings.Contains(*ver, "/") {
		x := strings.SplitN(*ver, "/", 2)
		if len(x) != 2 {
			fatal("Specify function/code separated by \"/\" ")
		}
		fn = x[0]
		cname = x[1]
	} else {
		cname = *ver
		fn = *vfn
	}

	var only_code bool
	only_code = *voc

	fnid := resolve(fcol, fn)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, cname)

	var ci api.CodeImage

	makeReq(xcol.Info(string(cver)), &ci)

	if !only_code {
		showInfoElement(elementCode{&ci})
	} else {
		fmt.Print(string(ci.Source.Text))
		fmt.Printf("\n")
	}
}

func parseCode(src string, ci *api.SourceImage) {
	if strings.HasPrefix(src, "http") {
		ci.URL = src
	} else if strings.HasPrefix(src, "repo:") {
		x := strings.SplitN(src, ":", 3)
		if len(x) != 3 {
			fatal("Specify repository as repo:<id>:<path>")
		}
		y := resolve(repcol, x[1])
		ci.RepoId = api.ObjectId(y)
		ci.Path = x[2]
	} else {
		var err error

		ci.Text, err = ioutil.ReadFile(src)
		if err != nil {
			fatal("Error reading sources: %s\n", err.Error())
		}
	}
}
