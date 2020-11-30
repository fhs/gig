// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm file ...",
	Short: "Remove files from the working tree and from the index",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	RunE:  gitRm,
}

func gitRm(_ *cobra.Command, args []string) error {
	root, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	for _, a := range args {
		a, err = repoRelPath(root, a)
		if err != nil {
			return err
		}
		_, err = w.Remove(a)
		if err != nil {
			return err
		}
	}
	return nil
}
