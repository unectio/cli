package main

import (
	"fmt"
	"flag"
	"strings"
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
)

var seccol = apilet.Secrets

func doSecret(cmd int, name *string) {
	sec_actions := map[int]func(*string) {}

	sec_actions[CmdAdd] = secretAdd
	sec_actions[CmdList] = secretList
	sec_actions[CmdInfo] = secretInfo
	sec_actions[CmdDel] = secretDelete

	doTargetCmd(cmd, name, sec_actions)
}

func parseKV(kv string) map[string]string {
	ret := make(map[string]string)

	for _, x := range strings.Split(kv, ";") {
		if x == "" {
			continue
		}
		y := strings.SplitN(x, "=", 2)
		ret[y[0]] = y[1]
	}

	return ret
}

func secretAdd(name *string) {
	kv := flag.String("kv", "", "table (k=v;...)")
	flag.Parse()

	pl := parseKV(*kv)

	sec := api.SecretImage {}

	sec.Name = generate(*name, "sec")
	sec.Payload = pl

	makeReq(seccol.Add(&sec), &sec)

	fmt.Printf("Added secret (id %s)\n", sec.Id)
}

func secretDelete(name *string) {
	flag.Parse()

	secid := resolve(seccol, *name)

	makeReq(seccol.Delete(string(secid)), nil)
}

func secretList(_ *string) {
	var secs []*api.SecretImage

	flag.Parse()

	makeReq(seccol.List(), &secs)

	for _, sec := range secs {
		fmt.Printf("%s: %s\n", sec.Id, sec.Name)
	}
}

func secretInfo(name *string) {
	flag.Parse()

	secid := resolve(seccol, *name)

	var sec api.SecretImage

	makeReq(seccol.Info(string(secid)), &sec)

	fmt.Printf("Id:             %s\n", sec.Id)
	fmt.Printf("Name:           %s\n", sec.Name)
	fmt.Printf("Tags:           %s\n", strings.Join(sec.Tags, ","))
}
