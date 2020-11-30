// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	fdiff "github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/diff"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:     "diff [commit]",
		Aliases: []string{"di"},
		Short:   "Show changes between working tree, index, commits, etc.",
		Long: `If a commit argument is given, compare working tree with that commit.
Otherwise, compare working tree with the index (staging area for the
next commit).`,
		Args: cobra.MaximumNArgs(1),
		RunE: diffCmd,
	}
	rootCmd.AddCommand(cmd)
}

func diffCmd(_ *cobra.Command, args []string) error {
	root, r, err := openRepo()
	if err != nil {
		return err
	}

	if len(args) == 1 {
		return diffWithCommit(os.Stdout, r, root, plumbing.Revision(args[0]))
	}
	return diffWithIndex(os.Stdout, r, root)
}

func diffWithIndex(w io.Writer, r *git.Repository, root string) error {
	idx, err := r.Storer.Index()
	if err != nil {
		return err
	}
	iter := &indexEntriesIter{
		idx: idx,
		r:   r,
		k:   0,
	}
	return worktreeDiff(w, iter, root)
}

func diffWithCommit(w io.Writer, r *git.Repository, root string, rev plumbing.Revision) error {
	hash, err := r.ResolveRevision(rev)
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(*hash)
	if err != nil {
		return err
	}
	tree, err := commit.Tree()
	if err != nil {
		return err
	}
	return worktreeDiff(w, tree.Files(), root)
}

type fileIter interface {
	Next() (*object.File, error)
}

type indexEntriesIter struct {
	idx *index.Index
	r   *git.Repository
	k   int
}

func (i *indexEntriesIter) Next() (*object.File, error) {
	entries := i.idx.Entries
	if i.k >= len(entries) {
		return nil, io.EOF
	}
	if i.k < 0 {
		return nil, fmt.Errorf("index %v out of range", i.k)
	}
	e := entries[i.k]
	b, err := i.r.BlobObject(e.Hash)
	if err != nil {
		return nil, err
	}
	i.k++
	return object.NewFile(e.Name, e.Mode, b), nil
}

func worktreeDiff(w io.Writer, iter fileIter, root string) error {
	var filePatches []fdiff.FilePatch
	for {
		file, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fromContent, err := file.Contents()
		if err != nil {
			return err
		}
		b, err := ioutil.ReadFile(filepath.Join(root, file.Name))
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			b = nil
		}
		toContent := string(b)
		if fromContent != toContent {
			fp, err := fileDiff(os.Stdout, file, fromContent, toContent)
			if err != nil {
				return err
			}
			filePatches = append(filePatches, fp)
		}
	}
	ue := fdiff.NewUnifiedEncoder(w, fdiff.DefaultContextLines)
	return ue.Encode(&gigPatch{
		message:     "",
		filePatches: filePatches,
	})
}

func fileDiff(w io.Writer, f *object.File, a, b string) (fdiff.FilePatch, error) {
	diffs := diff.Do(a, b)
	var chunks []fdiff.Chunk
	for _, d := range diffs {
		var op fdiff.Operation
		switch d.Type {
		case diffmatchpatch.DiffEqual:
			op = fdiff.Equal
		case diffmatchpatch.DiffDelete:
			op = fdiff.Delete
		case diffmatchpatch.DiffInsert:
			op = fdiff.Add
		}
		chunks = append(chunks, &gigChunk{content: d.Text, op: op})
	}

	isBinary, err := f.IsBinary()
	if err != nil {
		return nil, err
	}
	fp := &gigFilePatch{
		isBinary: isBinary,
		from: &gigFile{
			hash: f.Hash,
			mode: f.Mode,
			path: f.Name,
		},
		to: &gigFile{
			hash: f.Hash, // TODO
			mode: f.Mode, // TODO
			path: f.Name,
		},
		chunks: chunks,
	}
	return fp, nil
}

type gigPatch struct {
	message     string
	filePatches []fdiff.FilePatch
}

func (p *gigPatch) FilePatches() []fdiff.FilePatch { return p.filePatches }
func (p *gigPatch) Message() string                { return p.message }

type gigFilePatch struct {
	isBinary bool
	from, to *gigFile
	chunks   []fdiff.Chunk
}

func (fp *gigFilePatch) IsBinary() bool               { return fp.isBinary }
func (fp *gigFilePatch) Files() (from, to fdiff.File) { return fp.from, fp.to }
func (fp *gigFilePatch) Chunks() []fdiff.Chunk        { return fp.chunks }

type gigFile struct {
	hash plumbing.Hash
	mode filemode.FileMode
	path string
}

func (f *gigFile) Hash() plumbing.Hash     { return f.hash }
func (f *gigFile) Mode() filemode.FileMode { return f.mode }
func (f *gigFile) Path() string            { return f.path }

type gigChunk struct {
	content string
	op      fdiff.Operation
}

func (c *gigChunk) Content() string       { return c.content }
func (c *gigChunk) Type() fdiff.Operation { return c.op }
