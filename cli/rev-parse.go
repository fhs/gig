// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func init() {
	var rpc revParseCmd

	cmd := &cobra.Command{
		Use:   "rev-parse [revision ...]",
		Short: "Parse and resolve revision to corresponding hash",
		Long:  ``,
		RunE:  rpc.run,
	}
	rootCmd.AddCommand(cmd)
	cmd.Flags().BoolVar(&rpc.abbrevRef, "abbrev-ref", false, "No-op")
	cmd.Flags().BoolVar(&rpc.gitDir, "git-dir", false, "Show the path to the .git directory")
}

type revParseCmd struct {
	abbrevRef bool
	gitDir    bool
}

func (rpc *revParseCmd) run(_ *cobra.Command, args []string) error {
	root, r, err := openRepo()
	if err != nil {
		return err
	}
	if rpc.gitDir {
		fmt.Printf("%v\n", filepath.Join(root, ".git"))
		// don't return here: args may not be empty
	}
	for _, rev := range args {
		h, err := r.ResolveRevision(plumbing.Revision(rev))
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", h)
	}
	return nil
}
