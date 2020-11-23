// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package gig

import (
	"os"
	"testing"

	"github.com/fhs/gig/cli"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"gig": gig,
	}))
}

func TestGig(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Setup: func(env *testscript.Env) error {
			env.Vars = append(
				env.Vars,
				"GIT_AUTHOR_NAME=Test User",
				"GIT_AUTHOR_EMAIL=testuser@example.com",
			)
			return nil
		},
	})
}

func gig() int {
	cli.Execute()
	return 0
}
