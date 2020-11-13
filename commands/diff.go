// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/diff"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show changes between commits",
	Long:  ``,
	RunE:  gitDiff,
}

func gitDiff(_ *cobra.Command, args []string) error {
	root, r, err := openRepo()
	if err != nil {
		return err
	}
	ref, err := r.Head()
	if err != nil {
		if err.Error() == "reference not found" {
			// No diff if there is no commits yet.
			return nil
		}
		return err
	}
	h := ref.Hash()
	commit, err := r.CommitObject(h)
	if err != nil {
		return err
	}
	tree, err := commit.Tree()
	if err != nil {
		return err
	}
	return tree.Files().ForEach(func(file *object.File) error {
		reader, err := file.Reader()
		if err != nil {
			return err
		}
		old, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		cur, err := ioutil.ReadFile(filepath.Join(root, file.Name))
		if err != nil {
			return err
		}
		if olds, curs := string(old), string(cur); olds != curs {
			fmt.Println(file.Name, file.Mode)
			lineDiff(os.Stdout, olds, curs)
		}
		return nil
	})
}

func lineDiff(w io.Writer, a, b string) {
	// how to print this in unified diff format?
	diffs := diff.Do(a, b)
	for _, d := range diffs {
		lines := d.Text
		if lines[len(lines)-1] == '\n' {
			lines = lines[:len(lines)-1]
		}
		chr := rune(' ')
		switch d.Type {
		case diffmatchpatch.DiffDelete:
			chr = '-'
		case diffmatchpatch.DiffInsert:
			chr = '+'
		}
		lines = strings.Replace(lines, "\n", string([]rune{'\n', chr}), -1)
		fmt.Fprintf(w, "%c%s\n", chr, lines)
	}
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
