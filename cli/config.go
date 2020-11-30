// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

func init() {
	var cc configCmd

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Get repository or global options",
		Long:  `Gig uses the same configuration files as git(1).`,
		RunE:  cc.run,
	}
	rootCmd.AddCommand(cmd)
	cmd.Flags().BoolVarP(&cc.list, "list", "l", false, "List all variables set in config file")
}

type configCmd struct {
	list bool
}

func (cc *configCmd) run(_ *cobra.Command, args []string) error {
	if !cc.list {
		return fmt.Errorf("no flags given")
	}
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	cfg, err := r.ConfigScoped(config.GlobalScope)
	if err != nil {
		return err
	}
	fmt.Printf("user.name=%v\n", cfg.User.Name)
	fmt.Printf("user.email=%v\n", cfg.User.Email)
	return nil
}
