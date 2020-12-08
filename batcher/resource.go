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
	"github.com/unectio/util/sync"
)

const (
	debug = true
)

type resource interface {
	create(ApiletClient) error
	name() string
	id() api.ObjectId
}

type resWrap struct {
	res   resource
	ready chan struct{}
}

var resources = make(map[string]*resWrap)

func addResource(typ string, res resource) {
	rw := &resWrap{}
	rw.res = res
	rw.ready = make(chan struct{})

	resources[typ+"."+res.name()] = rw
}

func (rw *resWrap) add(cln ApiletClient, key string) error {
	err := rw.res.create(cln)
	if err != nil && rw.res.id() != "" {
		err = fmt.Errorf("misbehaved %s", key)
	}

	close(rw.ready)
	return err
}

func createResources(cln ApiletClient) error {
	var wg sync.WaitGroupErr

	for key, res := range resources {
		wg.Inc()
		go func(res *resWrap) {
			wg.Done(res.add(cln, key))
		}(res) /* copy the resource! */
	}

	return wg.Wait()
}

func waitResourceId(typ string, name_id *api.ObjectId, wg *sync.WaitGroupErr) {
	wg.Inc()
	go func() {
		var err error
		*name_id, err = waitId(typ, *name_id)
		wg.Done(err)
	}()
}

func waitId(typ string, name api.ObjectId) (api.ObjectId, error) {
	rw, ok := resources[typ+"."+string(name)]
	if !ok {
		return "", fmt.Errorf("cannot find %s %s", typ, string(name))
	}

	<-rw.ready

	id := rw.res.id()
	if id == "" {
		return "", fmt.Errorf("wait %s %s", typ, string(name))
	}

	return id, nil
}
