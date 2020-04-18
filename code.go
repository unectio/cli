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
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
	"io/ioutil"
	"strings"
)

var ccol = apilet.FnCodes

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

func codeAdd(fcname *string, vfn *string, vlang *string, vsrc *string, vw *int) {
	var fn, cname, lang, src string
	var w int

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

	fn = *fcname
	fid := resolve(fcol, fn)

	var cis []*api.CodeImage

	makeReq(ccol.Sub(string(fid)).List(), &cis)

	for _, ci := range cis {
		showListElement(elementCode{ci})
	}
}

func codeDel(fn *string, cname *string) {

	fnid := resolve(fcol, *fn)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, *cname)

	makeReq(xcol.Delete(string(cver)), nil)
}

func codeSet(fn *string, cn *string, vsrc *string, vw *int) {

	var ci api.CodeImage

	ci.Weight = *vw
	ci.Source = &api.SourceImage{}
	parseCode(*vsrc, ci.Source)

	fnid := resolve(fcol, *fn)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, *cn)

	makeReq(xcol.Upd(string(cver), &ci), nil)
}

func codeInfo(fn *string, cn *string, just_code *bool) {

	var only_code bool
	only_code = *just_code

	fnid := resolve(fcol, *fn)
	xcol := ccol.Sub(string(fnid))
	cver := resolve(xcol, *cn)

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
