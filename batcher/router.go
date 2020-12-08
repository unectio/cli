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
	"github.com/unectio/util/sync"
)

type specRtHandler struct{}

func (_ specRtHandler) parse(e *api.SpecEntry) (resource, error) {
	var rt api.RouterImage

	err := e.Spec.Unmarshal(&rt)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Collected router %s\n", rt.Name)

	return &router{&rt}, nil
}

type router struct {
	*api.RouterImage
}

func (rt *router) name() string     { return rt.Name }
func (rt *router) id() api.ObjectId { return rt.Id }

func (rt *router) create(cln ApiletClient) error {
	var wg sync.WaitGroupErr

	if rt.AuthId != "" {
		waitResourceId("auth", &rt.AuthId, &wg)
	}

	for _, rr := range rt.Mux {
		waitResourceId("function", &rr.FnId, &wg)
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	return makeAddReq(cln, apilet.Routers, &rt.RouterImage)
}
