// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone url directory",
	Short: "Clone a repository into a new directory",
	Long:  ``,
	RunE:  gitClone,
}

func gitClone(_ *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: gig clone <url> <dir>")
	}
	// Clone the given repository to the given directory
	_, err := git.PlainClone(args[1], false, &git.CloneOptions{
		URL:      args[0],
		Progress: os.Stdout,
	})
	return err
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
