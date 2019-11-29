package main

import (
	"fmt"
	"flag"
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
)

var authcol = apilet.AuthMethods

func doAuth(cmd int, name *string) {
	a_actions := map[int]func(namep *string) {}
	a_actions[CmdAdd] = authAdd
	a_actions[CmdList] = authList
	a_actions[CmdInfo] = authInfo
	a_actions[CmdDel] = authDel

	doTargetCmd(cmd, name, a_actions)
}

func authAdd(name *string) {
	key := flag.String("k", "", "jwt key (base64-encoded or auto)")
	flag.Parse()

	ai := api.AuthMethodImage{}
	ai.Name = generate(*name, "am")

	if *key != "" {
		ai.JWT = &api.AuthJWTImage { Key: *key }
	}

	makeReq(authcol.Add(&ai), &ai)

	fmt.Printf("Added auth method (id %s)\n", ai.Id)
}

func authDel(name *string) {
	flag.Parse()

	aid := resolve(authcol, *name)

	makeReq(authcol.Delete(string(aid)), nil)
}

func authList(_ *string) {
	var as []*api.AuthMethodImage

	makeReq(authcol.List(), &as)

	for _, am := range as {
		fmt.Printf("%s: %s\n", am.Id, am.Name)
	}
}

func authInfo(name *string) {
	flag.Parse()

	aid := resolve(authcol, *name)

	var ai api.AuthMethodImage

	makeReq(authcol.Info(string(aid)), &ai)

	fmt.Printf("Id:             %s\n", ai.Id)
	fmt.Printf("Name:           %s\n", ai.Name)
	if ai.JWT != nil {
		fmt.Printf("Key:            %s\n", ai.JWT.Key)
	}
}
