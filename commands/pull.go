// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package commands

import (
	git "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull [remote]",
	Short: "Fetch from and integrate with another repository",
	Long: `Incorporate changes from a remote repository into the current branch.
If remote is not specified, the remote named origin is used.`,
	Args: cobra.MaximumNArgs(1),
	RunE: gitPull,
}

func gitPull(cmd *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	remote := "origin"
	if len(args) >= 1 {
		remote = args[0]
	}
	return w.Pull(&git.PullOptions{
		RemoteName: remote,
	})
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
