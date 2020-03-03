//////////////////////////////////////////////////////////////////////////////
//
// (C) Copyright 2019-2020 by Unectio, Inc.
//
// The information contained herein is confidential, proprietary to Unectio,
// Inc.
//
//////////////////////////////////////////////////////////////////////////////

package main

import (
	"os"
	"fmt"
	"flag"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/unectio/api"
)

func main() {
	method := flag.String("m", "POST", "Method")
	path := flag.String("p", "", "URL path")
	body_from := flag.String("b", "", "Where to read body from (- means from stdin)")
	ctype := flag.String("t", "", "Content-type (works only with -b)")
	key := flag.String("k", "", "Trigger or route rule key")
	claims := flag.String("c", "", "Claims in 'name=value:...' format")
	flag.Parse()

	rq := api.Request{}
	rq.Method = *method
	rq.Path = *path
	rq.Key = *key

	if *body_from != "" {
		if *body_from == "-" {
			*body_from = "/dev/stdin"
		}

		data, err := ioutil.ReadFile(*body_from)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read body: %s\n", err.Error())
			os.Exit(1)
		}

		rq.Body = data
		rq.Content = *ctype
	}

	rq.Args = make(map[string]string)
	for _, arg := range flag.Args() {
		x := strings.SplitN(arg, "=", 2)
		if len(x) != 2 {
			fmt.Fprintf(os.Stderr, "Bad argument %s, expect foo=bar format\n", arg)
			os.Exit(1)
		}

		rq.Args[x[0]] = x[1]
	}

	if *claims != ""  {
		rq.Claims = make(map[string]interface{})
		for _, cl := range strings.Split(*claims, ":") {
			x := strings.SplitN(cl, "=", 2)
			if len(x) != 2 {
				fmt.Fprintf(os.Stderr, "Bad claim %s, expect foo=bar format\n", cl)
				os.Exit(1)
			}
			rq.Claims[x[0]] = x[1]
		}
	}

	result, _ := Main(&rq)

	res, err := json.Marshal(result)
	if err == nil {
		fmt.Printf("Returned: [%s]\n", res);
	}

}
