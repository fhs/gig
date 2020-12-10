// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"strconv"

	git "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "blame file",
		Short: "Show what revision and author last modified each line of a file",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE:  blameCmd,
	}
	rootCmd.AddCommand(cmd)
}

func blameCmd(_ *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	head, err := r.Head()
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(head.Hash())
	if err != nil {
		return err
	}
	result, err := git.Blame(commit, args[0])
	if err != nil {
		return err
	}
	pad := len(strconv.Itoa(len(result.Lines) - 1))
	for i, l := range result.Lines {
		fmt.Printf("%v (%v %v %*d) %v\n", l.Hash, l.Author, l.Date, pad, i, l.Text)
	}
	return nil
}
