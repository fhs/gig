// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

//go:build !plan9
// +build !plan9

package cli

import (
	"io"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func remoteAuth(r *git.Repository, remote string) (transport.AuthMethod, error) {
	return nil, nil // the default one works already
}

func endpointAuth(ep *transport.Endpoint) (transport.AuthMethod, error) {
	return nil, nil // the default one works already
}

var progressWriter = func() io.Writer {
	return os.Stdout
}
