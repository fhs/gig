// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

type commitCmd struct {
	message string
}

func (cc *commitCmd) run(_ *cobra.Command, args []string) error {
	root, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	status, err := w.Status()
	if err != nil {
		return err
	}
	if nothingInStaging(status) {
		return fmt.Errorf(`nothing to commit (use "gig add")`)
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
	if cc.message == "" {
		mfile := filepath.Join(root, ".git", "COMMIT_EDITMSG")
		err := ioutil.WriteFile(mfile, []byte(emptyCommit), 0644)
		if err != nil {
			return err
		}
		err = editFile(mfile)
		if err != nil {
			return fmt.Errorf("editor failed: %v", err)
		}
		msg, err := readCommit(mfile)
		if err != nil {
			return err
		}
		if msg == "" {
			return fmt.Errorf("aborting commit due to empty commit message")
		}
		cc.message = msg
	}
	_, err = w.Commit(cc.message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	})
	return err
}

func init() {
	var cc commitCmd

	cmd := &cobra.Command{
		Use:     "commit",
		Aliases: []string{"ci"},
		Short:   "Record changes to the repository",
		Long: `If the commit message (-m flag) is empty, a text editor will be opened to
edit the commit message. The text editor is configured by the either the
$GIT_EDITOR, $VISUAL, or $EDITOR environment variable, in that order,
falling back to a default editor if all three environment variables
are empty.
`,
		RunE: cc.run,
	}
	cmd.Flags().StringVarP(&cc.message, "message", "m", "", "Commit message")
	rootCmd.AddCommand(cmd)
}

func nothingInStaging(s git.Status) bool {
	for _, status := range s {
		switch status.Staging {
		case git.Unmodified, git.Untracked:
		default:
			return false
		}
	}
	return true
}

func editFile(filename string) error {
	args := strings.Fields(preferredEditor())
	if len(args) == 0 {
		panic("internal error: empty editor")
	}
	args = append(args, filename)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func readCommit(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var msg []byte
	for scanner.Scan() {
		b := append(scanner.Bytes(), '\n')
		if b[0] != '#' {
			msg = append(msg, b...)
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	msg = bytes.TrimFunc(msg, func(r rune) bool { return r == '\n' })
	if len(msg) > 0 {
		msg = append(msg, '\n')
	}
	return string(msg), nil
}

func preferredEditor() string {
	for _, name := range []string{
		"GIT_EDITOR",
		"VISUAL",
		"EDITOR",
	} {
		if e := os.Getenv(name); e != "" {
			return e
		}
	}
	return defaultEditor()
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

var emptyCommit = `
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
`
