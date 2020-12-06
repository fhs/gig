// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "push [[repository] refspec...]",
		Short: "Update remote refs along with associated objects",
		Long: `If not arguments are given, update the remote repository named origin.

If no refspec is given, pushes current branch.`,
		RunE: pushCmd,
	}
	rootCmd.AddCommand(cmd)
}

func pushCmd(cmd *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}

	remote := "origin"
	if len(args) >= 1 {
		remote = args[0]
	}

	var refspecs []config.RefSpec
	if len(args) > 1 {
		for _, s := range args[1:] {
			refspecs = append(refspecs, config.RefSpec(s))
		}
	} else {
		// No refspec was specified, so we want to push current branch.
		// Nil refspec seems to push all branches.
		head, err := r.Head()
		if err != nil {
			return err
		}
		refspecs = []config.RefSpec{
			config.RefSpec(head.Name() + ":" + head.Name()),
		}
	}
	auth, err := remoteAuth(r, remote)
	if err != nil {
		return err
	}
	err = r.Push(&git.PushOptions{
		RemoteName: remote,
		RefSpecs:   refspecs,
		Auth:       auth,
	})
	if err == git.NoErrAlreadyUpToDate {
		fmt.Printf("%v\n", err)
		return nil
	}
	return err
}
