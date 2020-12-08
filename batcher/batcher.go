//////////////////////////////////////////////////////////////////////////////
//
// (C) Copyright 2019-2020 by Unectio, Inc.
//
// The information contained herein is confidential, proprietary to Unectio,
// Inc.
//
//////////////////////////////////////////////////////////////////////////////

package batcher

import (
	"io"

	"github.com/unectio/api"
	"github.com/unectio/util"
)

func Process(in io.ReadCloser, cln ApiletClient) error {
	err := collectFile(in)
	if err != nil {
		return util.Error("Error parsing file", err)
	}

	return createResources(cln)
}

type specHandler interface {
	parse(*api.SpecEntry) (resource, error)
}

var specHandlers = map[string]specHandler{
	"function": specFnHandler{},
	"router":   specRtHandler{},
	"auth":     specAMHandler{},
}
