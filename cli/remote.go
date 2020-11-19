// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

type remoteCmd struct {
	verbose bool
}

func (rc *remoteCmd) list(_ *cobra.Command, args []string) error {
	_, repo, err := openRepo()
	if err != nil {
		return err
	}
	remotes, err := repo.Remotes()
	if err != nil {
		return err
	}
	for _, rem := range remotes {
		if rc.verbose {
			fmt.Printf("%v\n", rem)
		} else {
			cfg := rem.Config()
			fmt.Printf("%v\n", cfg.Name)
		}
	}
	return nil
}

func (rc *remoteCmd) add(_ *cobra.Command, args []string) error {
	_, repo, err := openRepo()
	if err != nil {
		return err
	}
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: args[0],
		URLs: []string{args[1]},
	})
	return err
}

func (rc *remoteCmd) remove(_ *cobra.Command, args []string) error {
	_, repo, err := openRepo()
	if err != nil {
		return err
	}
	return repo.DeleteRemote(args[0])
}

func init() {
	var rc remoteCmd

	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Manage set of tracked repositories",
		Long:  ``,
		RunE:  rc.list,
	}
	rootCmd.AddCommand(cmd)
	cmd.Flags().BoolVarP(&rc.verbose, "verbose", "v", false, "Be more verbose")

	cmd.AddCommand(&cobra.Command{
		Use:   "add name url",
		Short: "Add a remote named name for the repository at url",
		Long:  ``,
		Args:  cobra.ExactArgs(2),
		RunE:  rc.add,
	}, &cobra.Command{
		Use:   "remove name",
		Short: "Remove the remote named name",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE:  rc.remove,
	})
}
