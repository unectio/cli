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

type specAMHandler struct{}

func (_ specAMHandler) parse(e *api.SpecEntry) (resource, error) {
	var am api.AuthMethodImage

	err := e.Spec.Unmarshal(&am)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Collected auth method %s\n", am.Name)

	return &auth{&am}, nil
}

type auth struct {
	*api.AuthMethodImage
}

func (am *auth) name() string     { return am.Name }
func (am *auth) id() api.ObjectId { return am.Id }

func (am *auth) create(cln ApiletClient) error {
	return makeAddReq(cln, apilet.AuthMethods, &am.AuthMethodImage)
}
