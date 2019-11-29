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

	doTargetCmd(cmd, name, rep_actions)
}

func repoAdd(name *string) {
	url := flag.String("u", "", "repo URL (git)")
	flag.Parse()

	rp := api.RepoImage{}

	rp.Name = generate(*name, "repo")
	rp.Type = "git"
	rp.URL = *url

	makeReq(repcol.Add(&rp), &rp)

	fmt.Printf("Added repo (id %s)\n", rp.Id)
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
