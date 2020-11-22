// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

type logCmd struct {
	n int
}

func (lc *logCmd) run(_ *cobra.Command, args []string) error {
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
	iter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return err
	}
	defer iter.Close()
	for i := 0; lc.n < 0 || i < lc.n; i++ {
		c, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(c)
	}
	return nil
}

func init() {
	var lc logCmd

	cmd := &cobra.Command{
		Use:   "log",
		Short: "Show commit logs",
		Long:  ``,
		RunE:  lc.run,
	}
	rootCmd.AddCommand(cmd)
	cmd.Flags().IntVarP(&lc.n, "max-count", "n", -1, "Limit the number of commits to output")
}
