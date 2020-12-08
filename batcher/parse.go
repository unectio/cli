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
	"io"

	"github.com/unectio/api"
	"github.com/unectio/util"
	"gopkg.in/yaml.v2"
)

func collectFile(in io.ReadCloser) error {
	defer in.Close()

	for yd := range util.SplitYAML(in) {
		var entry api.SpecEntry

		err := yaml.Unmarshal(yd, &entry)
		if err != nil {
			return err
		}

		handler, ok := specHandlers[entry.Type]
		if !ok {
			return fmt.Errorf("No handler for %s type", entry.Type)
		}

		res, err := handler.parse(&entry)
		if err != nil {
			return err
		}

		addResource(entry.Type, res)
	}

	return nil
}
