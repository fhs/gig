// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Create an empty Git repository",
	Long:  ``,
	RunE:  gitInit,
}

func gitInit(_ *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: gig init [directory]")
	}
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	dotgit := filepath.Join(dir, ".git")
	fs := osfs.New(dotgit)
	s := filesystem.NewStorage(fs, nil)
	_, err := git.Init(s, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Initialized empty Git repository in %v\n", dotgit)
	return nil
}
