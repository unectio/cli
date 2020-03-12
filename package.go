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
	"net/url"
	"os"
)

var pcols = apilet.PkgLists

func doPackage(cmd int, name *string) {
	pk_actions := map[int]func(name *string){}

	pk_actions[CmdAdd] = packageAdd
	pk_actions[CmdDel] = packageDelete
	pk_actions[CmdList] = packageList
	pk_actions[CmdInfo] = packageInfo

	doTargetCmd(cmd, name, pk_actions)
}

type elementPk struct{ *api.PkgImage }

func (pe elementPk) id() string {
	return pe.PkgImage.Name
}

func (pe elementPk) name() string {
	return pe.PkgImage.Name
}

func (pe elementPk) version() string {
	return pe.PkgImage.Version
}

func (pe elementPk) short() string {
	return ""
}

func (pe elementPk) long() []*field {
	return []*field{
		{
			name: "Name",
			data: pe.PkgImage.Name,
		},
		{
			name: "Version",
			data: pe.PkgImage.Version,
		},
	}
}

func packageAdd(name *string) {
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	var lang = goopt.String([]string{"-l", "--language"}, "", "language of package")
	var ver = goopt.String([]string{"-v", "--version"}, "", "version of package")
	goopt.Parse(nil)

	if *lang == "" {
		fatal("Specify language")
	}

	pa := api.PkgImage{}
	pa.Name = *name
	pa.Version = *ver

	makeReq(pcols.Sub(*lang).Add(&pa), &pa)

	showAddedElement(elementPk{&pa})
}

func packageList(_ *string) {
	var pks []*api.PkgImage
	var lang = goopt.String([]string{"-l", "--language"}, "", "language of package")

	goopt.Summary = fmt.Sprintf("Usage: %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2])
	goopt.ExtraUsage = ""
	goopt.Parse(nil)

	makeReq(pcols.Sub(*lang).List(), &pks)

	for _, pk := range pks {
		showListElement(elementPk{pk})
	}
}

func packageInfo(name *string) {
	var pk api.PkgImage
	var lang = goopt.String([]string{"-l", "--language"}, "", "language of package")
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	goopt.Parse(nil)

	makeReq(pcols.Sub(*lang).Info(url.PathEscape(*name)), &pk)

	showInfoElement(elementPk{&pk})
}

func packageDelete(name *string) {
	var lang = goopt.String([]string{"-l", "--language"}, "", "language of package")

	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	goopt.Parse(nil)

	makeReq(pcols.Sub(*lang).Delete(*name), nil)
}
