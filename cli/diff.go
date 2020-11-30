// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	fdiff "github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/diff"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(diffCmd)
}

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:     "diff",
	Aliases: []string{"di"},
	Short:   "Show changes between commits",
	Long:    ``,
	RunE:    gitDiff,
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
	// TODO: We need to compare worktree with staging instead of HEAD.
	return worktreeDiff(os.Stdout, tree, root)
}

func worktreeDiff(w io.Writer, from *object.Tree, root string) error {
	iter := from.Files()
	defer iter.Close()
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
