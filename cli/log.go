// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

type logCmd struct {
	n      int
	format string
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
		if lc.format == "" {
			fmt.Println(c)
		} else {
			var b bytes.Buffer
			if err := formatCommit(&b, lc.format, c); err != nil {
				return err
			}
			fmt.Println(b.String())
		}
	}
	return nil
}

func formatCommit(w io.Writer, format string, c *object.Commit) error {
	i := strings.IndexByte(format, ':')
	if i < 0 || format[:i] != "format" {
		return fmt.Errorf("unsupported format string %q", format)
	}
	format = format[i+1:]

	fb := []byte(format)
	for len(fb) > 0 {
		i := bytes.IndexByte(fb, '%')
		if i < 0 || i == len(fb)-1 {
			_, err := w.Write(fb)
			return err
		}
		if i > 0 {
			_, err := w.Write(fb[:i])
			if err != nil {
				return err
			}
			fb = fb[i:]
		}
		fb = fb[1:] // skip '%'

		i = 0
		var err error
		switch fb[0] {
		case '%':
			_, err = fmt.Fprintf(w, "%%")
			i++

		case 'H', 'h': // commit hash
			// TODO: %h should be abbreviated commit hash, but
			// we're just trying to get `go tool dist` working for now.
			_, err = fmt.Fprintf(w, "%v", c.Hash)
			i++

		case 'c':
			if len(fb) > 1 {
				switch fb[1] {
				case 'd':
					_, err = fmt.Fprintf(w, "%v", c.Committer.When)
					i += 2
				}
			}
		}
		if i == 0 { // no expansion
			_, err = fmt.Fprintf(w, "%%")
		}
		if err != nil {
			return err
		}
		fb = fb[i:]
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
	cmd.Flags().StringVar(&lc.format, "format", "", "Print commits in the given format")
}
