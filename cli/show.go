// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [object]",
	Short: "Show various types of objects",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	RunE:  gitShow,
}

func gitShow(_ *cobra.Command, args []string) error {
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
	commit, err := r.CommitObject(*h)
	if err != nil {
		return err
	}
	fmt.Println(commit.String())

	// TODO: Implement patch for initial commit (0 parents)
	// and merge commits (2 parents).
	if commit.NumParents() == 1 {
		parent, err := commit.Parent(0)
		if err != nil {
			return err
		}
		patch, err := parent.Patch(commit)
		if err != nil {
			return err
		}
		fmt.Println(patch)
	}
	return nil
}
