// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "clone url directory",
		Short: "Clone a repository into a new directory",
		Long:  ``,
		RunE:  cloneCmd,
	}
	rootCmd.AddCommand(cmd)
}

func cloneCmd(_ *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: gig clone <url> <dir>")
	}
	auth, err := endpointAuth(args[0])
	if err != nil {
		return err
	}
	// Clone the given repository to the given directory
	_, err = git.PlainClone(args[1], false, &git.CloneOptions{
		URL:      args[0],
		Auth:     auth,
		Progress: progressWriter,
	})
	return err
}
