// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

// +build windows

package client

import (
	"context"
	"path/filepath"

	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol/keybase1"
)

func doSimpleFSPlatformGlob(g *libkb.GlobalContext, ctx context.Context, cli SimpleFSInterface, path keybase1.Path) ([]keybase1.Path, error) {
	var returnPaths []keybase1.Path

	// local glob
	matches, err := filepath.Glob(path.Local())
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		returnPaths = append(returnPaths, keybase1.NewPathWithLocal(match))
	}
	return returnPaths, nil
}
