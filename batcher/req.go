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
	rq "github.com/unectio/util/request"
	"github.com/unectio/util/restmux/client"
)

type ApiletClient interface {
	MakeRequest(*rq.Request, interface{}) error
}

func makeAddReq(cln ApiletClient, col *client.Collection, o interface{}) error {
	/*
	 * For now -- fire all requests in parallel
	 */
	return cln.MakeRequest(col.Add(o), o)
}
