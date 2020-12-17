// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "show [object]",
		Short: "Show various types of objects",
		Long: `Show an object, which can be a reference to a commit or a file within
a commit tree with syntax 'commit:filename'.`,
		Args: cobra.MaximumNArgs(1),
		RunE: showCmd,
	}
	rootCmd.AddCommand(cmd)
}

func showCmd(_ *cobra.Command, args []string) error {
	root, r, err := openRepo()
	if err != nil {
		return err
	}
	obj := "HEAD"
	if len(args) > 0 {
		obj = args[0]
	}
	filename := ""
	f := strings.SplitN(obj, ":", 2)
	if len(f) >= 2 {
		obj = f[0]
		filename = f[1]
	}

	h, err := r.ResolveRevision(plumbing.Revision(obj))
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(*h)
	if err != nil {
		return err
	}

	if filename != "" {
		if strings.HasPrefix(filename, "./") || strings.HasPrefix(filename, "../") {
			filename, err = repoRelPath(root, filename)
			if err != nil {
				return err
			}
		}
		f, err := commit.File(filename)
		if err != nil {
			return err
		}
		rd, err := f.Blob.Reader()
		if err != nil {
			return err
		}
		defer rd.Close()
		_, err = io.Copy(os.Stdout, rd)
		return err
	}

	fmt.Println(commit.String())

	// TODO: Implement patch for initial commit (0 parents)
	// and merge commits (2 parents).
	if commit.NumParents() == 1 {
		parent, err := commit.Parent(0)
		if err != nil {
			return err
		}
		patch, err := parent.Patch(commit)
		if err != nil {
			return err
		}
		fmt.Println(patch)
	}
	return nil
}
