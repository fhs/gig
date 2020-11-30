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
	var tc tagCmd

	cmd := &cobra.Command{
		Use:   "tag [name...]",
		Short: "Create, list, or delete a tag",
		Long: `If no arguments are given, list existing tags.

If a tag name is given, create that tag pointing to current HEAD.
`,
		RunE: tc.run,
	}
	rootCmd.AddCommand(cmd)

	cmd.Flags().BoolVarP(&tc.delete, "delete", "d", false, "Delete existing tags with the given names")
}

type tagCmd struct {
	delete bool
}

func (tc *tagCmd) run(_ *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}

	if tc.delete {
		for _, name := range args {
			if err := r.DeleteTag(name); err != nil {
				return err
			}
		}
		return nil
	}

	if len(args) == 0 { // list tags
		tagrefs, err := r.Tags()
		if err != nil {
			return err
		}
		return tagrefs.ForEach(func(t *plumbing.Reference) error {
			fmt.Println(t)
			return nil
		})
	}

	if got, want := len(args), 1; got != want {
		return fmt.Errorf("expected %v argument but got %v", want, got)
	}
	head, err := r.Head()
	if err != nil {
		return fmt.Errorf("HEAD: %v", err)
	}
	_, err = r.CreateTag(args[0], head.Hash(), nil)
	return err
}
