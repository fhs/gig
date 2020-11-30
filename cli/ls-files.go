// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/spf13/cobra"
)

type LsFilesCmd struct {
	stage bool
}

func (lfc *LsFilesCmd) run(_ *cobra.Command, args []string) error {
	root, err := findRepoRoot()
	if err != nil {
		return err
	}
	f, err := os.Open(filepath.Join(root, ".git", "index"))
	if err != nil {
		return err
	}
	defer f.Close()

	var idx index.Index
	dec := index.NewDecoder(f)
	if err := dec.Decode(&idx); err != nil {
		return err
	}
	for _, e := range idx.Entries {
		if lfc.stage {
			fmt.Printf("%06o %s %d\t%s\n", e.Mode, e.Hash, e.Stage, e.Name)
		} else {
			fmt.Printf("%v\n", e.Name)
		}
	}
	return nil
}

func init() {
	var lfc LsFilesCmd

	cmd := &cobra.Command{
		Use:   "ls-files",
		Short: "Show information about files in the index",
		Long:  ``,
		RunE:  lfc.run,
	}
	rootCmd.AddCommand(cmd)
	cmd.Flags().BoolVarP(&lfc.stage, "stage", "s", false, "Show staged contents' mode bits, object name and stage number")
}
