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
		log.Printf("-> %s\n", rq.Hdrs())
		if rq.Body != nil {
			x, _ := json.Marshal(rq.Body)
			log.Printf("BODY -> %s\n", string(x))
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
		if resb, ok := res.(*[]byte); ok {
			resp, *resb = resp.Raw()
		} else {
			resp = resp.B(res)
		}
	}

	if *debug {
		log.Printf("<- %s\n", resp.String())
		log.Printf("<- %s\n", resp.Hdrs())
		if res != nil {
			x, _ := json.Marshal(res)
			log.Printf("BODY <- %s\n", string(x))
		}
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
