// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"regexp"

	git "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func init() {
	var gc grepCmd

	cmd := &cobra.Command{
		Use:   "grep pattern",
		Short: "Print lines matching a pattern",
		Long:  `Search for the given regular expression pattern in the repository files.`,
		Args:  cobra.ExactArgs(1),
		RunE:  gc.run,
	}
	rootCmd.AddCommand(cmd)

	cmd.Flags().BoolVarP(&gc.lineNumber, "line-number", "n", false, "Print line number")
	cmd.Flags().BoolVarP(&gc.invertMatch, "invert-match", "v", false, "Select non-matching lines")
}

type grepCmd struct {
	lineNumber  bool
	invertMatch bool
}

func (gc *grepCmd) run(_ *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	reg, err := regexp.Compile(args[0])
	if err != nil {
		return err
	}
	results, err := w.Grep(&git.GrepOptions{
		Patterns:    []*regexp.Regexp{reg},
		InvertMatch: gc.invertMatch,
	})
	if err != nil {
		return err
	}
	for _, res := range results {
		if gc.lineNumber {
			fmt.Printf("%v:%v:%v\n", res.FileName, res.LineNumber, res.Content)
		} else {
			fmt.Printf("%v:%v\n", res.FileName, res.Content)
		}
	}
	return nil
}
