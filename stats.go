package main

import (
	"fmt"
	"flag"
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
)

var stprop = apilet.FnStats

func functionStats(name *string) {
	flag.Parse()

	fid := resolve(fcol, *name)

	var st api.FuncStatsImage

	makeReq(stprop.Get(string(fid)), &st)

	fmt.Printf("Calls:          %d\n", st.Calls)
	fmt.Printf("Run time:       %d us\n", st.RunTime)
	fmt.Printf("Last called:    %s\n", st.LastCall)
}
