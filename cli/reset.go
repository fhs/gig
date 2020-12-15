// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func init() {
	var rc resetCmd

	cmd := &cobra.Command{
		Use:   "reset [commit]",
		Short: "Reset current HEAD to the specified state",
		Long: `Set the current branch head to the specified commit (default HEAD),
and optionally modify the index and the working tree to match depending
on the mode: mixed (default), hard, merge or soft.`,
		Args: cobra.MaximumNArgs(1),
		RunE: rc.run,
	}
	rootCmd.AddCommand(cmd)

	cmd.Flags().BoolVar(&rc.mixed, "mixed", false, "Reset index but not working tree (default)")
	cmd.Flags().BoolVar(&rc.hard, "hard", false, "Reset index and working tree")
	cmd.Flags().BoolVar(&rc.merge, "merge", false, "Reset index but keep unmerged changes in working tree")
	cmd.Flags().BoolVar(&rc.soft, "soft", false, "Index and working tree are not changed")
}

type resetCmd struct {
	mixed, hard, merge, soft bool
}

func (gc *resetCmd) run(_ *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	obj := "HEAD"
	if len(args) > 0 {
		obj = args[0]
	}
	h, err := r.ResolveRevision(plumbing.Revision(obj))
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	mode, err := gc.mode()
	if err != nil {
		return err
	}
	return w.Reset(&git.ResetOptions{
		Commit: *h,
		Mode:   mode,
	})
}

func (gc *resetCmd) mode() (git.ResetMode, error) {
	var n int
	var m git.ResetMode
	if gc.mixed {
		m = git.MixedReset
		n++
	}
	if gc.hard {
		m = git.HardReset
		n++
	}
	if gc.merge {
		m = git.MergeReset
		n++
	}
	if gc.soft {
		m = git.SoftReset
		n++
	}
	if n == 0 {
		return git.MixedReset, nil
	}
	if n != 1 {
		return 0, fmt.Errorf("exactly one mode should be specified, not %v", n)
	}
	return m, nil
}
