// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"crypto"
	_ "crypto/sha256"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

const Date = "2006-01-02"

func init() {
	var pc pgpCmd

	cmd := &cobra.Command{
		Use:   "pgp",
		Short: "Manage OpenPGP keys",
		Long:  ``,
	}
	rootCmd.AddCommand(cmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "generate name comment email",
		Short: "Generate a key",
		Long:  ``,
		Args:  cobra.ExactArgs(3),
		RunE:  pc.generate,
	}, &cobra.Command{
		Use:   "list",
		Short: "List keys",
		Long:  ``,
		RunE:  pc.list,
	}, &cobra.Command{
		Use:   "list-secret",
		Short: "List secret keys",
		Long:  ``,
		RunE:  pc.listSecret,
	}, &cobra.Command{
		Use:   "read-armored",
		Short: "read armored key ring",
		Long:  ``,
		RunE:  pc.readArmored,
	}, &cobra.Command{
		Use:   "sign ring key",
		Short: "sign a message",
		Long:  ``,
		Args:  cobra.ExactArgs(2),
		RunE:  pc.sign,
	})
}

type pgpCmd struct {
}

func (pc *pgpCmd) generate(_ *cobra.Command, args []string) error {
	ent, err := openpgp.NewEntity(args[0], args[1], args[2], &packet.Config{
		DefaultHash: crypto.SHA256,
	})
	if err != nil {
		return err
	}
	err = ent.SerializePrivate(os.Stdout, nil)
	if err != nil {
		return err
	}
	return nil
}

func (pc *pgpCmd) list(_ *cobra.Command, args []string) error {
	ring, err := openpgp.ReadKeyRing(os.Stdin)
	if err != nil {
		return err
	}
	for _, ent := range ring {
		pk := ent.PrimaryKey
		nbits, err := pk.BitLength()
		if err != nil {
			return err
		}
		fmt.Printf("pub\t%T %v %v\n", pk.PublicKey, nbits, pk.CreationTime.Format(Date))
		fmt.Printf("\t%X\n", pk.Fingerprint)
		for key := range ent.Identities {
			fmt.Printf("uid\t%v\n", key)
		}
		for _, sub := range ent.Subkeys {
			pk := sub.PublicKey
			nbits, err := pk.BitLength()
			if err != nil {
				return err
			}
			fmt.Printf("sub\t%T %v %v\n", pk.PublicKey, nbits, pk.CreationTime.Format(Date))
		}
		fmt.Printf("\n")
	}
	return nil
}

func (pc *pgpCmd) listSecret(_ *cobra.Command, args []string) error {
	ring, err := openpgp.ReadKeyRing(os.Stdin)
	if err != nil {
		return err
	}
	for _, ent := range ring {
		sk := ent.PrivateKey
		nbits, err := sk.BitLength()
		if err != nil {
			return err
		}
		fmt.Printf("pub\t%T %v %v\n", sk.PrivateKey, nbits, sk.CreationTime.Format(Date))
		fmt.Printf("\t%X\n", sk.Fingerprint)
		for key := range ent.Identities {
			fmt.Printf("uid\t%v\n", key)
		}
		for _, sub := range ent.Subkeys {
			sk := sub.PrivateKey
			nbits, err := sk.BitLength()
			if err != nil {
				return err
			}
			fmt.Printf("sub\t%T %v %v\n", sk.PrivateKey, nbits, sk.CreationTime.Format(Date))
		}
		fmt.Printf("\n")
	}
	return nil
}

func (pc *pgpCmd) readArmored(_ *cobra.Command, args []string) error {
	ring, err := openpgp.ReadArmoredKeyRing(os.Stdin)
	if err != nil {
		return err
	}
	for _, ent := range ring {
		err = ent.Serialize(os.Stdout)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pc *pgpCmd) sign(_ *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer f.Close()

	ring, err := openpgp.ReadKeyRing(f)
	if err != nil {
		return err
	}
	id, err := strconv.ParseUint(args[1], 16, 64)
	if err != nil {
		return err
	}
	keys := ring.KeysByIdUsage(id, packet.KeyFlagSign)
	if len(keys) != 1 {
		return fmt.Errorf("expected 1 key but found %v", len(keys))
	}
	err = openpgp.ArmoredDetachSign(os.Stdout, keys[0].Entity, os.Stdin, nil)
	if err != nil {
		return err
	}
	return nil
}
