// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"os"
	"sort"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "ls-remote [remote]",
		Short: "List references in a remote repository",
		Long: `List references in a remote repository along with the associted commit
hash. If the remote argument is not provided, origin is used.`,
		Args: cobra.MaximumNArgs(1),
		RunE: lsRemoteCmd,
	}
	rootCmd.AddCommand(cmd)
}

func lsRemoteCmd(cmd *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	name := "origin"
	if len(args) >= 1 {
		name = args[0]
	}
	rem, err := r.Remote(name)
	if err != nil {
		return err
	}
	auth, err := remoteAuth(r, name)
	if err != nil {
		return err
	}
	refs, err := rem.List(&git.ListOptions{
		Auth: auth,
	})
	if err != nil {
		return err
	}

	if cfg := rem.Config(); len(cfg.URLs) > 0 {
		fmt.Fprintf(os.Stderr, "From %v\n", cfg.URLs[0])
	}
	sort.Slice(refs, func(i, j int) bool {
		return refs[i].Name() < refs[j].Name()
	})
	m := newRefMap(refs)
	for _, ref := range refs {
		fmt.Printf("%v\t%v\n", m.resolveHash(ref), ref.Name())
	}
	return nil
}

type refMap map[plumbing.ReferenceName]*plumbing.Reference

func newRefMap(refs []*plumbing.Reference) refMap {
	m := make(refMap, len(refs))
	for _, ref := range refs {
		m[ref.Name()] = ref
	}
	return m
}

func (m refMap) resolveHash(ref *plumbing.Reference) plumbing.Hash {
	for ref.Type() == plumbing.SymbolicReference {
		var ok bool
		ref, ok = m[ref.Target()]
		if !ok {
			return plumbing.ZeroHash
		}
	}
	return ref.Hash()
}
