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

var rfcols = apilet.RepoFiles

func doRepoFile(cmd int, name *string) {
	fl_actions := map[int]func(name *string){}

	fl_actions[CmdList] = fileList
	fl_actions[CmdInfo] = fileInfo

	doTargetCmd(cmd, name, fl_actions)
}

type elementFl struct{ *api.RepoFileImage }

func (fe elementFl) id() string {
	return string(fe.RepoFileImage.Name)
}

func (fe elementFl) name() string {
	return fe.RepoFileImage.Name
}

func (fe elementFl) version() string {
	return fe.RepoFileImage.Name
}

func (fe elementFl) short() string {
	return fe.RepoFileImage.Path
}

func (fe elementFl) long() []*field {
	return []*field{
		{
			name: "Path",
			data: fe.RepoFileImage.Path,
		},
		{
			name: "Type",
			data: fe.RepoFileImage.Type,
		},
	}
}

func fileList(_ *string) {
	var fs []*api.RepoFileImage
	var repo_id = goopt.String([]string{"-r", "--repository"}, "", "repository")

	goopt.Summary = fmt.Sprintf("Usage: %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2])
	goopt.ExtraUsage = ""
	goopt.Parse(nil)

	rpid := resolve(repcol, *repo_id)
	makeReq(rfcols.Sub(string(rpid)).List(), &fs)

	for _, ff := range fs {
		fmt.Printf("%s\n", ff.Path)
	}
}

func fileInfo(fname *string) {
	var rf api.RepoFileImage
	goopt.Summary = fmt.Sprintf("Usage: %s %s %s %s:\n", os.Args[0], os.Args[1], os.Args[2], os.Args[3])
	goopt.ExtraUsage = ""
	var voc = goopt.Flag([]string{"-C", "--code"}, []string{}, "Show code only", "")
	var repo_id = goopt.String([]string{"-r", "--repository"}, "", "repository")

	goopt.Parse(nil)

	var only_code bool
	only_code = *voc

	rpid := resolve(repcol, *repo_id)
	makeReq(rfcols.Sub(string(rpid)).Info(*fname), &rf)

	if !only_code {
		//showInfoElement(elementCode{&rf})
	} else {
		//fmt.Print(string(rf.Source.Text))
		fmt.Printf("\n")
	}
}
