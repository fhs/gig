// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"bufio"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fhs/go-plan9-auth/auth"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	tssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func remoteAuth(r *git.Repository, remote string) (transport.AuthMethod, error) {
	rem, err := r.Remote(remote)
	if err != nil {
		return nil, err
	}
	cfg := rem.Config()
	if len(cfg.URLs) == 0 {
		return nil, fmt.Errorf("not URLs in remote %v", remote)
	}
	return endpointAuth(cfg.URLs[0])
}

func endpointAuth(url string) (transport.AuthMethod, error) {
	t, err := transport.NewEndpoint(url)
	if err != nil {
		return nil, err
	}
	if t.Protocol != "ssh" {
		// We only care about ssh auth for now.
		// In the future, factotum should be used for https also.
		return nil, nil
	}
	// Plan 9 ssh client stores known host keys in the format
	// described by thumbprint(6).	An initial attempt to reuse this
	// file proved not very useful because the Plan 9 ssh client seems
	// to be using a different algorithm than golang.org/x/crypto/ssh
	// (ssh-rsa instead of ecdsa-sha2-nistp256), leading to thumbprint
	// mismatch. We also need to write less code if we just use the
	// OpenSSH known_hosts format.
	files, err := getDefaultKnownHostsFiles()
	if err != nil {
		return nil, err
	}
	hostKeyCallback, err := knownhosts.New(files...)
	if err != nil {
		return nil, err
	}
	return &tssh.PublicKeysCallback{
		User:     t.User,
		Callback: factotumSigners,
		HostKeyCallbackHelper: tssh.HostKeyCallbackHelper{
			HostKeyCallback: printUnknownKey(hostKeyCallback),
		},
	}, nil
}

func getDefaultKnownHostsFiles() ([]string, error) {
	files := filepath.SplitList(os.Getenv("SSH_KNOWN_HOSTS"))
	if len(files) != 0 {
		return files, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return []string{
		filepath.Join(home, "lib", "ssh", "known_hosts"),
	}, nil
}

func printUnknownKey(f ssh.HostKeyCallback) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		err := f(hostname, remote, key)
		if ke, ok := err.(*knownhosts.KeyError); ok && len(ke.Want) == 0 {
			b := base64.StdEncoding.EncodeToString(key.Marshal())
			host := hostname
			if r := remote.String(); r != host {
				host += "," + r
			}
			fmt.Fprintf(os.Stderr, "Add host key to known_hosts file after verification:\n")
			fmt.Fprintf(os.Stderr, "\techo '%v %v %v' >> $home/lib/ssh/known_hosts\n", host, key.Type(), b)
		}
		return err
	}
}

func factotumSigners() ([]ssh.Signer, error) {
	f, err := os.Open("/mnt/factotum/ctl")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var signers []ssh.Signer
	for scanner.Scan() {
		m, err := parseKey(scanner.Text())
		if err != nil {
			return nil, err
		}
		if m["proto"] == "rsa" && m["service"] == "ssh" {
			n, ok := new(big.Int).SetString(m["n"], 16)
			if !ok {
				return nil, fmt.Errorf("could not parse ek in rsa key")
			}
			ek, err := strconv.ParseUint(m["ek"], 16, 0)
			if err != nil {
				return nil, err
			}
			sn, err := newRSASigner(n, int(ek))
			if err != nil {
				return nil, err
			}
			signers = append(signers, sn)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return signers, nil
}

func parseKey(s string) (map[string]string, error) {
	f := strings.Fields(s)
	if len(f) == 0 || f[0] != "key" {
		return nil, fmt.Errorf("invalid factotum key")
	}
	m := make(map[string]string, len(f)-1)
	for _, e := range f {
		if i := strings.IndexByte(e, '='); i >= 0 {
			m[e[:i]] = e[i+1:]
		}
	}
	return m, nil
}

type rsaSigner struct {
	pub  *rsa.PublicKey
	spub ssh.PublicKey
}

func newRSASigner(N *big.Int, E int) (ssh.Signer, error) {
	pub := &rsa.PublicKey{
		N: N, // n in factotum
		E: E, // ek in factotum
	}
	spub, err := ssh.NewPublicKey(pub)
	if err != nil {
		return nil, err
	}
	return &rsaSigner{
		pub:  pub,
		spub: spub,
	}, nil
}

func (s *rsaSigner) PublicKey() ssh.PublicKey {
	return s.spub
}

func (s *rsaSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	rpc, err := auth.NewRPC()
	if err != nil {
		return nil, err
	}
	a := fmt.Sprintf("proto=rsa role=sign n=%X ek=%X", s.pub.N, s.pub.E)
	if _, _, err = rpc.Call("start", []byte(a)); err != nil {
		return nil, err
	}
	digest := sha1.Sum(data)
	if _, _, err = rpc.Call("write", digest[:]); err != nil {
		return nil, err
	}
	_, blob, err := rpc.Call("read", nil)
	if err != nil {
		return nil, err
	}
	return &ssh.Signature{
		Format: "ssh-rsa",
		Blob:   blob,
	}, nil
}
