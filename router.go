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
	"os"
	"fmt"
	"flag"
	"bufio"
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
	rt_actions[CmdUpdate] = routerUpdate
	rt_actions[CmdDel] = routerDelete

	doTargetCmd(cmd, name, rt_actions)
}

func formatRule(rule *api.RouteRuleImage, fnames map[api.ObjectId]string) string {
	return fmt.Sprintf("\t%s/%s=%s\n", rule.Methods, rule.Path, fnames[rule.FnId])
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

/*
 * Table can be provided either as a option argument (plain string)
 * or be written in a file.
 *
 * The plain string is comma-separated set of rules, the file is
 * one rule by line.
 *
 * Rule is M/P=F where M is the comma-separated list of methods,
 * the P is slash-separated path and F is function name.
 *
 * If you want to edit a mux on a server, there's no such API. Instead,
 * do this:
 *
 *  $ show router -M > mux.txt
 *  $ edit it with your favourite text editor
 *  $ upd router -tf mux.txt
 *
 * or this:
 *
 *  $ show router -M | sed -e ... | upd router -tf -
 *
 */
func parseTable(table, file string) ([]*api.RouteRuleImage, error) {
	var ret []*api.RouteRuleImage

	if table != "" {
		for _, rule := range strings.Split(table, ":") {
			r, err := parseRule(rule)
			if err != nil {
				return nil, err
			}

			ret = append(ret, r)
		}
	} else if file != "" {
		var f *os.File

		if file == "-" {
			f = os.Stdin
		} else {
			var err error

			f, err = os.Open(file)
			if err != nil {
				return nil, errors.New("error reading table file: " + err.Error())
			}

			defer f.Close()
		}

		sc := bufio.NewScanner(f)
		for sc.Scan() {
			r, err := parseRule(sc.Text())
			if err != nil {
				return nil, err
			}

			ret = append(ret, r)
		}
	} else {
		return nil, errors.New("either table or file needed")
	}

	return ret, nil
}

var muxprop = apilet.RtMux

func routerUpdate(name *string) {
	table := flag.String("t", "", "table (m,.../path=fn:...)")
	table_from := flag.String("tf", "", "file to read table from (in info -M format)")
	flag.Parse()

	rtid := resolve(rtcol, *name)

	if *table != "" || *table_from != "" {
		mux, err := parseTable(*table, *table_from)
		if err != nil {
			fatal("Error parsing table: %s\n", err.Error())
		}

		makeReq(muxprop.Set(string(rtid), mux), nil)
	}
}

func routerAdd(name *string) {
	table := flag.String("t", "", "table (m,.../path=fn:...)")
	table_from := flag.String("tf", "", "file to read table from (in info -M format)")
	url := flag.String("u", "", "custom URL to work on")
	flag.Parse()

	mux, err := parseTable(*table, *table_from)
	if err != nil {
		fatal("Error parsing table: %s\n", err.Error())
	}

	rt := api.RouterImage{}

	rt.Name = generate(*name, "router")
	rt.Mux = mux

	if *url != "" {
		rt.URL = api.URLProjectPfx + *url
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
	mux_only := flag.Bool("M", false, "Show only the mux")
	flag.Parse()

	rtid := resolve(rtcol, *name)

	var rt api.RouterImage

	makeReq(rtcol.Info(string(rtid)), &rt)

	if !*mux_only {
		fmt.Printf("Id:             %s\n", rt.Id)
		fmt.Printf("Name:           %s\n", rt.Name)
		fmt.Printf("URL:            %s\n", rt.URL)
		fmt.Printf("Table:\n")
	}

	var fns []*api.FunctionImage
	makeReq(fcol.List(), &fns)

	fnames := make(map[api.ObjectId]string)
	for _, fn := range fns {
		fnames[fn.Id] = fn.Name
	}

	for _, rule := range rt.Mux {
		fmt.Printf("\t%s\n", formatRule(rule, fnames))
	}
}
