// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import "os"

func defaultEditor() string {
	_, err := os.Stat("/dev/wsys")
	if err != nil && os.IsNotExist(err) { // console
		return "ed"
	}
	return "acme -c 1"
}
