// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add file ...",
	Short: "Add file contents to the index",
	Long:  ``,
	RunE:  gitAdd,
}

func gitAdd(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		fmt.Printf("Nothing specified, nothing added.")
		return nil
	}

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
		_, err = w.Add(a)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
