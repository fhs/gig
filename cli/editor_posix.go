// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

//go:build !plan9
// +build !plan9

package cli

import "os/exec"

func defaultEditor() string {
	for _, e := range []string{
		"vim",
		"vi",
		"nano",
	} {
		if p, err := exec.LookPath(e); err == nil {
			return p
		}
	}
	return "ed"
}
