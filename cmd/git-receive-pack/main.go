// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/fhs/gig/cli"
)

func main() {
	os.Args = append([]string{"gig", "receive-pack"}, os.Args[1:]...)
	cli.Execute()
}
