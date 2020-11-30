// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st"},
	Short:   "Show the working tree status",
	Long:    ``,
	RunE:    gitStatus,
}

func gitStatus(_ *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	// TODO: this is the --porcelain output
	// We also need to show paths relative to current working directory,
	// not relative to workdir root.
	status, err := w.Status()
	if err != nil {
		return err
	}
	fmt.Print(status)
	return nil
}
