// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing/transport/file"
	"github.com/spf13/cobra"
)

// uploadPackCmd represents the upload-pack command
var uploadPackCmd = &cobra.Command{
	Use:   "upload-pack directory",
	Short: "Send objects packed back to git-fetch-pack",
	Long:  ``,
	RunE:  gitUploadPack,
}

func gitUploadPack(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: upload-pack <dir>")
	}
	return file.ServeUploadPack(dotGitDir(args[0]))
}

func init() {
	rootCmd.AddCommand(uploadPackCmd)
}
