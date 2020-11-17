// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push [[repository] refspec...]",
	Short: "Update remote refs along with associated objects",
	Long: `If not arguments are given, update the remote repository named origin.

If no refspec is given, pushes current branch.`,
	RunE: gitPush,
}

func gitPush(cmd *cobra.Command, args []string) error {
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
	}
	fmt.Printf("refspecs = %v\n", refspecs)
	return r.Push(&git.PushOptions{
		RemoteName: remote,
		RefSpecs:   refspecs,
	})
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
