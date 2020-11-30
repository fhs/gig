// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "mv source destination",
		Short: "Move or rename a file or directory",
		Long:  ``,
		Args:  cobra.ExactArgs(2),
		RunE:  mvCmd,
	}
	rootCmd.AddCommand(cmd)
}

func mvCmd(_ *cobra.Command, args []string) error {
	root, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	src, err := repoRelPath(root, args[0])
	if err != nil {
		return err
	}
	dst, err := repoRelPath(root, args[1])
	if err != nil {
		return err
	}
	_, err = w.Move(src, dst)
	return err
}
