// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List branches",
	Long:  ``,
	RunE:  gitBranch,
}

func gitBranch(_ *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	bIter, err := r.Branches()
	if err != nil {
		return err
	}
	return bIter.ForEach(func(ref *plumbing.Reference) error {
		fmt.Println(ref)
		return nil
	})
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
