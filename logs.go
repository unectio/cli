package main

import (
	"fmt"
	"flag"
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
)

var logprop = apilet.FnLogs

func functionLogs(name *string) {
	lfor := flag.String("f", "", "for what period (duration since now)")
	flag.Parse()

	fid := resolve(fcol, *name)

	var logs []*api.LogEntry

	rq := logprop.Get(string(fid))
	if *lfor != "" {
		rq.Path += "?for=" + *lfor
	}

	makeReq(rq, &logs)

	for _, le := range logs {
		fmt.Printf("%s%16s: %s\n", le.Time, le.Event, le.Text)
	}
}
