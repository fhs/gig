// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing/transport/file"
	"github.com/spf13/cobra"
)

// receivePackCmd represents the receive-pack command
var receivePackCmd = &cobra.Command{
	Use:   "receive-pack directory",
	Short: "Receive what is pushed into the repository",
	Long:  ``,
	RunE:  gitReceivePack,
}

func gitReceivePack(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: receive-pack <dir>")
	}
	return file.ServeReceivePack(dotGitDir(args[0]))
}

func init() {
	rootCmd.AddCommand(receivePackCmd)
}
