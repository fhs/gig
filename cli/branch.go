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
	var bc branchCmd

	cmd := &cobra.Command{
		Use:     "branch [name...]",
		Aliases: []string{"br"},
		Short:   "List, create, or delete branches",
		Long: `With no arguments, list existing branches. The current branch is
prefixed with an asterisk.

If one argument is given, create a new branch named name which points
to the current HEAD.
`,
		RunE: bc.run,
	}
	rootCmd.AddCommand(cmd)

	cmd.Flags().BoolVarP(&bc.forceDelete, "force-delete", "D", false, "Force delete a branch")
}

type branchCmd struct {
	forceDelete bool // -D
}

func (bc *branchCmd) run(_ *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	head, err := r.Head()
	if err == plumbing.ErrReferenceNotFound {
		// maybe there are no commits yet
		head = nil
	} else if err != nil {
		return fmt.Errorf("HEAD: %v", err)
	}

	if bc.forceDelete {
		if len(args) == 0 {
			return fmt.Errorf("branch name required")
		}
		for _, name := range args {
			br := plumbing.NewBranchReferenceName(name)
			if head != nil && head.Name() == br {
				return fmt.Errorf("cannot delete checked out branch %q", br)
			}
			if err := r.Storer.RemoveReference(br); err != nil {
				return err
			}
		}
		return nil
	}

	if len(args) == 1 {
		if head == nil {
			return fmt.Errorf("reference to HEAD not found")
		}
		return r.Storer.SetReference(plumbing.NewHashReference(
			plumbing.NewBranchReferenceName(args[0]),
			head.Hash(),
		))
	}

	if len(args) != 0 {
		return fmt.Errorf("accepts 0 args, received %v", len(args))
	}
	bIter, err := r.Branches()
	if err != nil {
		return err
	}
	return bIter.ForEach(func(ref *plumbing.Reference) error {
		if head != nil && ref.Name() == head.Name() {
			fmt.Printf("* %v\n", ref)
		} else {
			fmt.Printf("  %v\n", ref)
		}
		return nil
	})
}
