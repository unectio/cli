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
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
)

var pcols = apilet.PkgLists

type elementPk struct{ *api.PkgImage }

func (pe elementPk) id() string {
	return string(pe.PkgImage.Id)
}

func (pe elementPk) name() string {
	return pe.PkgImage.Name
}

func (pe elementPk) version() string {
	return pe.PkgImage.Version
}

func (pe elementPk) short() string {
	return ""
}

func (pe elementPk) long() []*field {
	return []*field{
		{
			name: "Version",
			data: pe.PkgImage.Version,
		},
	}
}

func packageAdd(name *string, lang *string, ver *string) {
	pa := api.PkgImage{}
	pa.Name = *name
	pa.Version = *ver

	makeReq(pcols.Sub(*lang).Add(&pa), &pa)

	showAddedElement(elementPk{&pa})
}

func packageList(lang *string) {
	var pks []*api.PkgImage

	makeReq(pcols.Sub(*lang).List(), &pks)

	for _, pk := range pks {
		showListElement(elementPk{pk})
	}
}

func packageInfo(name *string, lang *string) {
	var pk api.PkgImage

	makeReq(pcols.Sub(*lang).Info(*name), &pk)

	showInfoElement(elementPk{&pk})
}

func packageDel(name *string, lang *string) {

	makeReq(pcols.Sub(*lang).Delete(*name), nil)
}
