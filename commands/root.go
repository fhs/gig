// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "gig",
	Short:        "A git clone implemented in Go",
	Long:         ``,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if !hasDeps() {
		os.Exit(2)
	}
	if err := rootCmd.Execute(); err != nil {
		// The error seems to be already printed by cobra
		os.Exit(1)
	}
}

func hasDeps() bool {
	ok := true
	for _, prog := range []string{
		transport.UploadPackServiceName,
		transport.ReceivePackServiceName,
	} {
		if _, err := exec.LookPath(prog); err != nil {
			fmt.Fprintf(os.Stderr, "%v is not installed. Install it from github.com/fhs/gig/cmd/%v.\n", prog, prog)
			ok = false
		}
	}
	return ok
}
