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
)

var authcol = apilet.AuthMethods

func authAdd(name *string, key *string) {

	ai := api.AuthMethodImage{}
	ai.Name = generate(*name, "am")

	if *key != "" {
		ai.JWT = &api.AuthJWTImage{Key: *key}
	}

	makeReq(authcol.Add(&ai), &ai)

	fmt.Printf("Added auth method (id %s)\n", ai.Id)
}

func authDel(name *string) {

	aid := resolve(authcol, *name)

	makeReq(authcol.Delete(string(aid)), nil)
}

func authList() {
	var as []*api.AuthMethodImage

	makeReq(authcol.List(), &as)

	for _, am := range as {
		fmt.Printf("%s: %s\n", am.Id, am.Name)
	}
}

func authInfo(name *string) {

	aid := resolve(authcol, *name)

	var ai api.AuthMethodImage

	makeReq(authcol.Info(string(aid)), &ai)

	fmt.Printf("Id:             %s\n", ai.Id)
	fmt.Printf("Name:           %s\n", ai.Name)
	if ai.JWT != nil {
		fmt.Printf("Key:            %s\n", ai.JWT.Key)
	}
}
