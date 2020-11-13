// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

func findRepoRoot() (string, error) {
	p, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		_, err := os.Stat(filepath.Join(p, ".git"))
		if err == nil {
			return p, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		newp := filepath.Join(p, "..")
		if newp == p {
			break
		}
		p = newp
	}
	return "", fmt.Errorf("fatal: not a git repository")
}

// dotGitDir returns the .git directory within dir.
// If there is no .git directory within dir, it returns dir.
func dotGitDir(dir string) string {
	if filepath.Base(dir) == ".git" {
		return dir
	}
	dot := filepath.Join(dir, ".git")
	fi, err := os.Stat(dot)
	if err == nil && fi.IsDir() {
		return dot
	}
	return dir
}

func openRepo() (string, *git.Repository, error) {
	root, err := findRepoRoot()
	if err != nil {
		return "", nil, err
	}
	r, err := git.Open(
		filesystem.NewStorage(
			osfs.New(filepath.Join(root, ".git")),
			cache.NewObjectLRUDefault(),
		),
		osfs.New(root),
	)
	return root, r, err
}
