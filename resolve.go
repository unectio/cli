package main

import (
	"github.com/unectio/api"
	"github.com/unectio/util/restmux/client"
)

func resolve(col *client.Collection, val string) api.ObjectId {
	if val == "" {
		fatal("No name to resolve given")
	}

	if val[0] == byte('%') {
		return api.ObjectId(val[1:])
	}

	var resp struct { Id api.ObjectId `json:"id"` }

	l, err := getLogin()
	if err != nil {
		fatal("Cannot login: %s", err.Error())
	}

	err = l.MakeRequest(col.Lookup(val), &resp)
	if err != nil {
		fatal("Cannot resolve %s/%s: %s", col, val, err.Error())
	}

	return resp.Id
}
