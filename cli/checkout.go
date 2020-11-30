// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func init() {
	var cc checkoutCmd

	cmd := &cobra.Command{
		Use:     "checkout branch",
		Aliases: []string{"co"},
		Short:   "Switch branches",
		Long:    ``,
		Args:    cobra.ExactArgs(1),
		RunE:    cc.run,
	}
	rootCmd.AddCommand(cmd)

	cmd.Flags().BoolVarP(&cc.create, "create", "b", false, "Create branch before checkout")
}

type checkoutCmd struct {
	create bool
}

func (cc *checkoutCmd) run(cmd *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	return w.Checkout(&git.CheckoutOptions{
		Hash:   plumbing.ZeroHash,
		Branch: plumbing.NewBranchReferenceName(args[0]),
		Create: cc.create,
		Force:  false,
		Keep:   true,
	})
}
