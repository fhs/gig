// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var lfc lsFilesCmd

	cmd := &cobra.Command{
		Use:   "ls-files",
		Short: "Show information about files in the index",
		Long:  ``,
		RunE:  lfc.run,
	}
	rootCmd.AddCommand(cmd)
	cmd.Flags().BoolVarP(&lfc.stage, "stage", "s", false, "Show staged contents' mode bits, object name and stage number")
}

type lsFilesCmd struct {
	stage bool
}

func (lfc *lsFilesCmd) run(_ *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	idx, err := r.Storer.Index()
	if err != nil {
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
