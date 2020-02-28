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
	}
}

func codeAdd(cname *string) {
	fn := flag.String("f", "", "function name/id")
	lang := flag.String("l", "", "language")
	src := flag.String("s", "", "sources (e.g. -- a file name)")
	flag.Parse()

	fid := resolve(fcol, *fn)

	var ci api.CodeImage

	ci.Name = generate(*cname, "code")
	ci.Lang = *lang
	ci.Source = &api.SourceImage{}
	parseCode(*src, ci.Source)

	makeReq(ccol.Sub(string(fid)).Add(&ci), &ci)

	showAddedElement(elementCode{&ci})
}

func codeList(_ *string) {
	fn := flag.String("f", "", "function name/id")
	flag.Parse()

	fid := resolve(fcol, *fn)

	var cis []*api.CodeImage

	makeReq(ccol.Sub(string(fid)).List(), &cis)

	for _, ci := range cis {
		showListElement(elementCode{ci})
	}
}

func codeDel(ver *string) {
	fn := flag.String("f", "", "function name/id")
	flag.Parse()

	fnid := resolve(fcol, *fn)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, *ver)

	makeReq(xcol.Delete(string(cver)), nil)
}

func codeUpdate(ver *string) {
	fn := flag.String("f", "", "function name/id")
	src := flag.String("s", "", "sources (e.g. -- a file name)")
	flag.Parse()

	var ci api.CodeImage

	ci.Source = &api.SourceImage{}
	parseCode(*src, ci.Source)

	fnid := resolve(fcol, *fn)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, *ver)

	makeReq(xcol.Upd(string(cver), &ci), nil)
}

func codeInfo(ver *string) {
	fn := flag.String("f", "", "function name/id")
	only_code := flag.Bool("C", false, "show code only")
	flag.Parse()

	fnid := resolve(fcol, *fn)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, *ver)

	var ci api.CodeImage

	makeReq(xcol.Info(string(cver)), &ci)

	if !*only_code {
		showInfoElement(elementCode{&ci})
	} else {
		fmt.Print(string(ci.Source.Text))
		fmt.Printf("\n")
	}
}
