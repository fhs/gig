package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/src-d/go-git/utils/diff"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
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

func gitAdd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Nothing specified, nothing added.")
	}

	root, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	for _, a := range args {
		a, err = filepath.Abs(a)
		if err != nil {
			return err
		}
		a, err = filepath.Rel(root, a)
		if err != nil {
			return err
		}
		_, err = w.Add(a)
		if err != nil {
			return err
		}
	}
	return nil
}

func gitBranch(args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	bIter, err := r.Branches()
	if err != nil {
		return err
	}
	return bIter.ForEach(func(ref *plumbing.Reference) error {
		fmt.Println(ref)
		return nil
	})
}

func gitClone(args []string) error {
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

func gitCommit(args []string) error {
	cl := flag.NewFlagSet("commit", flag.ExitOnError)
	msg := cl.String("m", "no message", "commit message")
	err := cl.Parse(args)
	if err != nil {
		cl.PrintDefaults()
		return err
	}

	_, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	u, err := user.Current()
	if err != nil {
		return err
	}
	_, err = w.Commit(*msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  u.Name,
			Email: u.Username + "@localhost.localdomain",
			When:  time.Now(),
		},
	})
	return err
}

func gitDiff(args []string) error {
	root, r, err := openRepo()
	if err != nil {
		return err
	}
	ref, err := r.Head()
	if err != nil {
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

func gitInit(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: gig init <name>")
	}
	fs := osfs.New(filepath.Join(args[0], ".git"))
	s := filesystem.NewStorage(fs, nil)
	_, err := git.Init(s, nil)
	return err
}

func gitLog(args []string) error {
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
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return err
	}
	return cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})
}

func gitShow(args []string) error {
	// TODO: support hash sub-string prefix
	// TODO: initial commit patch

	_, r, err := openRepo()
	if err != nil {
		return err
	}
	var h plumbing.Hash
	if len(args) < 1 {
		ref, err := r.Head()
		if err != nil {
			return err
		}
		h = ref.Hash()
	} else {
		h = plumbing.NewHash(args[0])
	}
	commit, err := r.CommitObject(h)
	if err != nil {
		return err
	}
	fmt.Println(commit.String())
	parent, err := commit.Parent(0)
	if err != nil {
		return err
	}
	patch, err := parent.Patch(commit)
	if err != nil {
		return err
	}
	fmt.Println(patch)
	return nil
}

func gitStatus(args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	// TODO: this is the --porcelain output
	status, err := w.Status()
	if err != nil {
		return err
	}
	fmt.Print(status)
	return nil
}

func run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no command given")
	}
	cmd := args[0]
	args = args[1:]

	switch cmd {
	case "add":
		return gitAdd(args)
	case "branch":
		return gitBranch(args)
	case "clone":
		return gitClone(args)
	case "commit":
		return gitCommit(args)
	case "diff":
		return gitDiff(args)
	case "init":
		return gitInit(args)
	case "log":
		return gitLog(args)
	case "show":
		return gitShow(args)
	case "status":
		return gitStatus(args)
	}
	return fmt.Errorf("unknown command %q", cmd)
}

func main() {
	flag.Parse()

	if err := run(flag.Args()); err != nil {
		log.Fatalln(err)
	}
}
