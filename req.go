package main

import (
	"log"
	"reflect"
	"github.com/unectio/api"
	"github.com/unectio/util"
	"encoding/json"
	"github.com/unectio/util/request"
)

func makeReq(rq *rq.Request, res interface{}) {
	l, err := getLogin()
	if err != nil {
		fatal("Cannot login: %s", err.Error())
	}

	err = l.MakeRequest(rq, res)
	if err != nil {
		fatal("Cannot make %s req: %s", rq.Path, err.Error())
	}
}

func (l *Login)MakeRequest(rq *rq.Request, res interface{}) error {
	rq.Host = l.address
	rq.Path = "/v1/" + rq.Path

	rq = rq.H(api.AuthTokHeader, util.BearerPrefix + l.token)
	if *debug {
		log.Printf("-> %s\n", rq.String())
		log.Printf("`- %s\n", rq.Hdrs())
		if rq.Body != nil {
			x, _ := json.Marshal(rq.Body)
			log.Printf("`- [%s]\n", string(x))
		}
	}

	if *dryrun {
		if rq.Method == "POST" {
			trySetSomeId(res)
		}

		return nil
	}

	resp := rq.Do()
	if res != nil {
		resp = resp.B(res)
	}

	if *debug {
		log.Printf("<- %s\n", resp.String())
	}

	return resp.E()
}

func trySetSomeId(x interface{}) {
	v := reflect.ValueOf(x)
	for v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	f := v.FieldByName("Id")
	if f.CanSet() {
		f.SetString("DEADBEEF")
	}
}
