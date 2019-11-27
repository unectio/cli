package main

import (
	"fmt"
	"flag"
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
)

var tgcol = apilet.FnTriggers

func doTrigger(cmd int, name *string) {
	tg_actions := map[int]func(*string) {}

	tg_actions[CmdAdd] = triggerAdd
	tg_actions[CmdList] = triggerList
	tg_actions[CmdInfo] = triggerInfo
	tg_actions[CmdDel] = triggerDel

	doTargetCmd(cmd, name, tg_actions)
}

type elementTg struct { *api.FuncTriggerImage }


func (te elementTg)id() string {
	return string(te.FuncTriggerImage.Id)
}

func (te elementTg)name() string {
	return te.FuncTriggerImage.Name
}

func (te elementTg)short() string {
	if te.FuncTriggerImage.URL != nil {
		return fmt.Sprintf("url: %s", te.FuncTriggerImage.URL.URL)
	}

	return ""
}

func (te elementTg)long() []*field {
	if te.FuncTriggerImage.URL != nil {
		return []*field {
			{
				name:	"URL",
				data:	te.FuncTriggerImage.URL.URL,
			},
			{
				name:	"Auth",
				data:	string(te.FuncTriggerImage.URL.AuthId),
			},
		}
	}

	return nil
}

func triggerAdd(name *string) {
	fn := flag.String("f", "", "function name/id")
	src := flag.String("s", "", "trigger source")
	auth := flag.String("a", "", "URL trigger auth name/id")
	flag.Parse()

	fid := resolve(fcol, *fn)

	tra := api.FuncTriggerImage{}
	tra.Name = *name

	switch *src {
	case "url":
		tra.URL = &api.URLTrigImage{ URL: "auto" }
		if *auth != "" {
			tra.URL.AuthId = resolve(authcol, *auth)
		}
	}

	makeReq(tgcol.Sub(string(fid)).Add(&tra), &tra)

	showAddedElement(elementTg{&tra})
}

func triggerList(_ *string) {
	fn := flag.String("f", "", "function name/id")
	flag.Parse()

	fid := resolve(fcol, *fn)

	var tgs []*api.FuncTriggerImage

	makeReq(tgcol.Sub(string(fid)).List(), &tgs)

	for _, tg := range tgs {
		showListElement(elementTg{tg})
	}
}

func triggerDel(name *string) {
	fn := flag.String("f", "", "function name/id")
	flag.Parse()

	fnid := resolve(fcol, *fn)
	tcol := tgcol.Sub(string(fnid))
	tgid := resolve(tcol, *name)

	makeReq(tcol.Delete(string(tgid)), nil)
}

func triggerInfo(name *string) {
	fn := flag.String("f", "", "function name/id")
	flag.Parse()

	fnid := resolve(fcol, *fn)
	tcol := tgcol.Sub(string(fnid))
	tgid := resolve(tcol, *name)

	var tg api.FuncTriggerImage

	makeReq(tcol.Info(string(tgid)), &tg)

	showInfoElement(elementTg{&tg})
}
