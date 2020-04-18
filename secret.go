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
	"fmt"
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
	"strings"
)

var seccol = apilet.Secrets

func parseKV(kv string) map[string]string {
	ret := make(map[string]string)

	for _, x := range strings.Split(kv, ",") {
		if x == "" {
			continue
		}
		y := strings.SplitN(x, "=", 2)
		ret[y[0]] = y[1]
	}

	return ret
}

func secretAdd(name *string, kv *string) {

	pl := parseKV(*kv)

	sec := api.SecretImage{}

	sec.Name = generate(*name, "sec")
	sec.Payload = pl
	sec.Reveal = "keys"

	makeReq(seccol.Add(&sec), &sec)

	fmt.Printf("Added secret (id %s)\n", sec.Id)
}

func secretDel(name *string) {

	secid := resolve(seccol, *name)

	makeReq(seccol.Delete(string(secid)), nil)
}

func secretList() {
	var secs []*api.SecretImage

	makeReq(seccol.List(), &secs)

	for _, sec := range secs {
		fmt.Printf("%s: %s\n", sec.Id, sec.Name)
	}
}

func secretInfo(name *string) {

	secid := resolve(seccol, *name)

	var sec api.SecretImage

	makeReq(seccol.Info(string(secid)), &sec)

	fmt.Printf("Id:             %s\n", sec.Id)
	fmt.Printf("Name:           %s\n", sec.Name)
	fmt.Printf("Tags:           %s\n", strings.Join(sec.Tags, ","))
	if sec.Payload != nil {
		fmt.Printf("------------------\nData:\n%-20s %-8s %-32s\n", "Key", "Value", "Reference")
		for key, _ := range sec.Payload {
			fmt.Printf("%-20s %-8s %-32s\n", key, "***", "$ref:secret."+sec.Name+"."+key)
		}
	}
}
