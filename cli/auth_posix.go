// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

// +build !plan9

package cli

import (
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func remoteAuth(r *git.Repository, remote string) (transport.AuthMethod, error) {
	return nil, nil // the default one works already
}

func endpointAuth(url string) (transport.AuthMethod, error) {
	return nil, nil // the default one works already
}
