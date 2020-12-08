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
	"fmt"

	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
)

type specFnHandler struct{}

func (_ specFnHandler) parse(e *api.SpecEntry) (resource, error) {
	var fn api.FunctionImage

	err := e.Spec.Unmarshal(&fn)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Collected function %s\n", fn.Name)

	return &function{&fn}, nil
}

type function struct {
	*api.FunctionImage
}

func (fn *function) name() string     { return fn.Name }
func (fn *function) id() api.ObjectId { return fn.Id }

func (fn *function) create(cln ApiletClient) error {
	return makeAddReq(cln, apilet.Functions, fn)
}
