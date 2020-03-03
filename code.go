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
	"io/ioutil"
	"strings"
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
)

var ccol = apilet.FnCodes

func doCode(cmd int, name *string) {
	co_actions := map[int]func(*string) {}

	co_actions[CmdAdd] = codeAdd
	co_actions[CmdList] = codeList
	co_actions[CmdDel] = codeDel
	co_actions[CmdUpdate] = codeUpdate
	co_actions[CmdInfo] = codeInfo

	doTargetCmd(cmd, name, co_actions)
}

type elementCode struct { *api.CodeImage }

func (ce elementCode)id() string {
	return string(ce.CodeImage.Id)
}

func (ce elementCode)name() string {
	return ce.CodeImage.Name
}

func (ce elementCode)short() string {
	return ce.CodeImage.Lang
}

func (ce elementCode)long() []*field {
	return []*field {
		{
			name:	"State",
			data:	ce.CodeImage.State,
		},
		{
			name:	"Gen",
			data:	ce.CodeImage.Gen,
		},
		{
			name:	"Lang",
			data:	ce.CodeImage.Lang,
		},
		{
			name:	"Weight",
			data:	ce.CodeImage.Weight,
		},
	}
}

func codeAdd(cname *string) {
	var fn string
	var lang string
	var src string
	var w int 
	const (
		fndefault_value = ""
		fnusage   = "function name/id"
		langdefault_value = ""
		langusage   = "language"
		srcdefault_value = ""
		srcusage   = "sources (file name or url or repo:<repo name>:path)"
		wdefault_value = 0
		wusage   = "code weight"
	)
	flag.StringVar(&fn, "function", fndefault_value, fnusage)
	flag.StringVar(&fn, "f", fndefault_value, fnusage+" (shorthand)")
	flag.StringVar(&lang, "language", langdefault_value, langusage)
	flag.StringVar(&lang, "l", langdefault_value, langusage+" (shorthand)")
	flag.StringVar(&src, "source", srcdefault_value, srcusage)
	flag.StringVar(&src, "s", srcdefault_value, srcusage+" (shorthand)")
	flag.IntVar(&w, "weight", wdefault_value, wusage)
	flag.IntVar(&w, "w", wdefault_value, wusage+" (shorthand)")
	flag.Parse()

	fid := resolve(fcol, fn)

	var ci api.CodeImage

	ci.Name = generate(*cname, "code")
	ci.Lang = lang
	ci.Weight = w
	ci.Source = &api.SourceImage{}
	parseCode(src, ci.Source)

	makeReq(ccol.Sub(string(fid)).Add(&ci), &ci)

	showAddedElement(elementCode{&ci})
}

func codeList(_ *string) {
	var fn string
	const (
		default_value = ""
		usage   = "function name/id"
	)
	flag.StringVar(&fn, "function", default_value, usage)
	flag.StringVar(&fn, "f", default_value, usage+" (shorthand)")
	flag.Parse()

	fid := resolve(fcol, fn)

	var cis []*api.CodeImage

	makeReq(ccol.Sub(string(fid)).List(), &cis)

	for _, ci := range cis {
		showListElement(elementCode{ci})
	}
}

func codeDel(ver *string) {
	var fn string
	const (
		default_value = ""
		usage   = "function name/id"
	)
	flag.StringVar(&fn, "function", default_value, usage)
	flag.StringVar(&fn, "f", default_value, usage+" (shorthand)")
	flag.Parse()

	fnid := resolve(fcol, fn)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, *ver)

	makeReq(xcol.Delete(string(cver)), nil)
}

func codeUpdate(ver *string) {
	var fn string
	var src string
	var w int 
	const (
		fndefault_value = ""
		fnusage   = "function name/id"
		srcdefault_value = ""
		srcusage   = "sources (file name or url or repo:<repo name>:path)"
		wdefault_value = 0
		wusage   = "code weight"
	)
	flag.StringVar(&fn, "function", fndefault_value, fnusage)
	flag.StringVar(&fn, "f", fndefault_value, fnusage+" (shorthand)")
	flag.StringVar(&src, "source", srcdefault_value, srcusage)
	flag.StringVar(&src, "s", srcdefault_value, srcusage+" (shorthand)")
	flag.IntVar(&w, "weight", wdefault_value, wusage)
	flag.IntVar(&w, "w", wdefault_value, wusage+" (shorthand)")
	flag.Parse()

	var ci api.CodeImage

	ci.Weight = w
	ci.Source = &api.SourceImage{}
	parseCode(src, ci.Source)

	fnid := resolve(fcol, fn)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, *ver)

	makeReq(xcol.Upd(string(cver), &ci), nil)
}

func codeInfo(ver *string) {
	var fn string
	var only_code bool
	const (
		fndefault_value = ""
		fnusage   = "function name/id"
		only_default_value = false
		only_usage   = "show code only"
	)
	flag.StringVar(&fn, "function", fndefault_value, fnusage)
	flag.StringVar(&fn, "f", fndefault_value, fnusage+" (shorthand)")
	flag.BoolVar(&only_code, "code", only_default_value, only_usage)
	flag.BoolVar(&only_code, "C", only_default_value, only_usage+" (shorthand)")
	flag.Parse()

	fnid := resolve(fcol, fn)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, *ver)

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

