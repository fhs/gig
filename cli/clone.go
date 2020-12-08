// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "clone url [directory]",
		Short: "Clone a repository into a new directory",
		Long: `Clone a repository into a new directory. If the directory is not privided,
it's derived from the URL.`,
		Args: cobra.MaximumNArgs(2),
		RunE: cloneCmd,
	}
	rootCmd.AddCommand(cmd)
}

func cloneCmd(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no URL provided")
	}
	url := args[0]

	ep, err := transport.NewEndpoint(url)
	if err != nil {
		return err
	}

	dir := ""
	if len(args) >= 2 {
		dir = args[1]
	} else {
		dir = cloneDir(ep.Path)
	}

	auth, err := endpointAuth(ep)
	if err != nil {
		return err
	}
	fmt.Printf("Cloning into %q...\n", dir)
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:      url,
		Auth:     auth,
		Progress: progressWriter,
	})
	return err
}

// cloneDir return the directory where newly cloned repository will be stored
// based on the full path in git URL. For example, it converts
// fhs/gig to gig, and fhs/gig.git to gig.
func cloneDir(dir string) string {
	d := path.Base(dir)
	if suf := ".git"; strings.HasSuffix(d, suf) && len(d) > len(suf) {
		d = d[:len(d)-len(suf)]
	}
	return d
}
