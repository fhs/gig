// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Long:  ``,
	RunE:  gitLog,
}

func gitLog(_ *cobra.Command, args []string) error {
	p, err := findRepoRoot()
	if err != nil {
		return err
	}
	r, err := git.PlainOpen(filepath.Join(p, ".git"))
	if err != nil {
		return err
	}
	ref, err := r.Head()
	if err != nil {
		return err
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return err
	}
	return cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})
}

func init() {
	rootCmd.AddCommand(logCmd)
}
