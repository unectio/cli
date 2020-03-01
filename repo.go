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

var repcol = apilet.Repos

func doRepo(cmd int, name *string) {
	rep_actions := map[int]func(*string) {}

	rep_actions[CmdAdd] = repoAdd
	rep_actions[CmdList] = repoList
	rep_actions[CmdInfo] = repoInfo
	rep_actions[CmdDel] = repoDel

	doTargetCmd(cmd, name, rep_actions)
}

func repoAdd(name *string) {
	url := flag.String("u", "", "repo URL (git)")
	flag.Parse()

	if *url == "" {
		fatal("No URL specified, mind using -u option")
	}

	rp := api.RepoImage{}

	rp.Name = generate(*name, "repo")
	rp.Type = "git"
	rp.URL = *url

	makeReq(repcol.Add(&rp), &rp)

	fmt.Printf("Added repo (id %s)\n", rp.Id)
}

func repoDel(name *string) {
	flag.Parse()

	rpid := resolve(repcol, *name)
	
	makeReq(repcol.Delete(string(rpid)), nil)
	
	fmt.Printf("Deleted repo (id %s)\n", rpid)
}

func repoList(_ *string) {
	var rps []*api.RepoImage

	flag.Parse()

	makeReq(repcol.List(), &rps)

	for _, rp := range rps {
		sfx := ""
		if rp.State != "ready" {
			sfx = " (" + rp.State + ")"
		}

		fmt.Printf("%s: %6s %s%s\n", rp.Id, rp.Type, rp.Name, sfx)
	}
}

func repoInfo(name *string) {
	flag.Parse()

	rpid := resolve(repcol, *name)

	var rp api.RepoImage

	makeReq(repcol.Info(string(rpid)), &rp)

	fmt.Printf("Id:             %s\n", rp.Id)
	fmt.Printf("Name:           %s\n", rp.Name)
	fmt.Printf("Type:           %s\n", rp.Type)
	fmt.Printf("URL:            %s\n", rp.URL)
	fmt.Printf("Head:           %s\n", rp.Head)
}
