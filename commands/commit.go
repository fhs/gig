// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package commands

import (
	"os/user"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

var message string

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:     "commit",
	Aliases: []string{"ci"},
	Short:   "Record changes to the repository",
	Long:    ``,
	RunE:    gitCommit,
}

func gitCommit(_ *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	u, err := user.Current()
	if err != nil {
		return err
	}
	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  u.Name,
			Email: u.Username + "@localhost.localdomain",
			When:  time.Now(),
		},
	})
	return err
}

func init() {
	commitCmd.Flags().StringVarP(&message, "message", "m", "no message", "Commit message")
	rootCmd.AddCommand(commitCmd)
}
