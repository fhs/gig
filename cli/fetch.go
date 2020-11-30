// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	git "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "fetch [remote]",
		Short: "Download objects and refs from another repository",
		Long:  `If remote is not specified, the remote named origin is used.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  fetchCmd,
	}
	rootCmd.AddCommand(cmd)
}

func fetchCmd(cmd *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	remote := "origin"
	if len(args) >= 1 {
		remote = args[0]
	}
	return r.Fetch(&git.FetchOptions{
		RemoteName: remote,
	})
}
