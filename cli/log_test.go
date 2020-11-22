// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"bytes"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestFormatCommit(t *testing.T) {
	for _, tc := range []struct {
		format string
		commit object.Commit
		want   string
	}{
		{
			"format: +%h %cd HEAD",
			object.Commit{
				Hash: plumbing.NewHash("9567f7269a075b80e01f4611402344b924a40c35"),
				Committer: object.Signature{
					Name:  "",
					Email: "",
					When:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				},
			},
			" +9567f7269a075b80e01f4611402344b924a40c35 2009-11-10 23:00:00 +0000 UTC HEAD",
		},
	} {
		var b bytes.Buffer
		if err := formatCommit(&b, tc.format, &tc.commit); err != nil {
			t.Fatal(err)
		}
		if got, want := b.String(), tc.want; got != want {
			t.Errorf("formatting %q: got %v, want %v", tc.format, got, want)
		}
	}
}
