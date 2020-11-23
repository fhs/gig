// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

var message string

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:     "commit",
	Aliases: []string{"ci"},
	Short:   "Record changes to the repository",
	Long:    ``,
	RunE:    gitCommit,
}

func gitCommit(_ *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	name := os.Getenv("GIT_AUTHOR_NAME")
	email := os.Getenv("GIT_AUTHOR_EMAIL")
	if name == "" || email == "" {
		cfg, err := r.ConfigScoped(config.GlobalScope)
		if err != nil {
			return err
		}
		if name == "" {
			name = cfg.User.Name
		}
		if email == "" {
			email = cfg.User.Email
		}
	}
	if name == "" || email == "" {
		fmt.Fprintf(os.Stderr, "%v\n", unknownUserMsg)
		return fmt.Errorf("user's name and/or email are empty")
	}
	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	})
	return err
}

func init() {
	commitCmd.Flags().StringVarP(&message, "message", "m", "no message", "Commit message")
	rootCmd.AddCommand(commitCmd)
}

var unknownUserMsg = `Author identity unknown

*** Please tell me who you are.

To set your account's default identity, write a configuration file like
the following to .git/config file or the global <HOME>/.gitconfig file:

	[user]
	email = gitster@example.com
	name = Junio C Hamano

Alternatively, set the GIT_AUTHOR_NAME and GIT_AUTHOR_EMAIL environment
variables instead.
`
