package main

import (
	"fmt"
	"flag"
	"errors"
	"strings"
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
)

var rtcol = apilet.Routers

func init() {
}

func doRouter(cmd int, name *string) {
	rt_actions := map[int]func(*string) {}

	rt_actions[CmdAdd] = routerAdd
	rt_actions[CmdList] = routerList
	rt_actions[CmdInfo] = routerInfo
	rt_actions[CmdDel] = routerDelete

	doTargetCmd(cmd, name, rt_actions)
}

func parseRule(rule string) (*api.RouteRuleImage, error) {
	sep1 := strings.Index(rule, "/")
	if sep1 == -1 {
		return nil, errors.New("No methods/path split")
	}

	sep2 := strings.LastIndex(rule, "=")
	if sep2 == -1 {
		return nil, errors.New("No path/function split")
	}

	fnid := resolve(fcol, rule[sep2+1:])

	ret := &api.RouteRuleImage {}
	ret.Methods = rule[:sep1]
	ret.Path = rule[sep1+1:sep2]
	ret.FnId = fnid

	return ret, nil
}

func parseTable(table string) ([]*api.RouteRuleImage, error) {
	var ret []*api.RouteRuleImage

	for _, rule := range strings.Split(table, ":") {
		r, err := parseRule(rule)
		if err != nil {
			return nil, err
		}

		ret = append(ret, r)
	}

	return ret, nil
}

func routerAdd(name *string) {
	table := flag.String("t", "", "table (m,.../path=fn:...)")
	url := flag.String("u", "", "custom URL to work on")
	flag.Parse()

	mux, err := parseTable(*table)
	if err != nil {
		fatal("Error parsing table: %s\n", err.Error())
	}

	rt := api.RouterImage{}

	rt.Name = generate(*name, "router")
	rt.Mux = mux

	if *url != "" {
		rt.URL = *url
	} else {
		rt.URL = api.AutoValue
	}

	makeReq(rtcol.Add(&rt), &rt)

	fmt.Printf("Added router (id %s)\n", rt.Id)
}

func routerDelete(name *string) {
	flag.Parse()

	rtid := resolve(rtcol, *name)

	makeReq(rtcol.Delete(string(rtid)), nil)
}

func routerList(_ *string) {
	var rts []*api.RouterImage

	flag.Parse()

	makeReq(rtcol.List(), &rts)

	for _, rt := range rts {
		fmt.Printf("%s: %s\n", rt.Id, rt.Name)
	}
}

func routerInfo(name *string) {
	flag.Parse()

	rtid := resolve(rtcol, *name)

	var rt api.RouterImage

	makeReq(rtcol.Info(string(rtid)), &rt)

	fmt.Printf("Id:             %s\n", rt.Id)
	fmt.Printf("Name:           %s\n", rt.Name)
	fmt.Printf("URL:            %s\n", rt.URL)
	fmt.Printf("Table:\n")
	for _, rule := range rt.Mux {
		fmt.Printf("\t%s/%s=%s\n", rule.Methods, rule.Path, rule.FnId)
	}
}
