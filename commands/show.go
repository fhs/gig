// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show various types of objects",
	Long:  ``,
	RunE:  gitShow,
}

func gitShow(_ *cobra.Command, args []string) error {
	// TODO: support hash sub-string prefix
	// TODO: initial commit patch

	_, r, err := openRepo()
	if err != nil {
		return err
	}
	var h plumbing.Hash
	if len(args) < 1 {
		ref, err := r.Head()
		if err != nil {
			return err
		}
		h = ref.Hash()
	} else {
		h = plumbing.NewHash(args[0])
	}
	commit, err := r.CommitObject(h)
	if err != nil {
		return err
	}
	fmt.Println(commit.String())
	parent, err := commit.Parent(0)
	if err != nil {
		return err
	}
	patch, err := parent.Patch(commit)
	if err != nil {
		return err
	}
	fmt.Println(patch)
	return nil
}

func init() {
	rootCmd.AddCommand(showCmd)
}
