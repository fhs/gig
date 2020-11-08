package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

var remotePort = flag.Int("p", 22, "Port to connect to on the remote host.")

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: %v [options] host\n", os.Args[0])
		os.Exit(2)
	}
	addr := net.JoinHostPort(args[0], strconv.Itoa(*remotePort))

	user, err := user.Current()
	if err != nil {
		log.Fatalf("could not get current user: %v\n", err)
	}

	agentClient, err := newAgentClient()
	if err != nil {
		log.Fatalf("could not connect to ssh-agent: %v", err)
	}

	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
	config := &ssh.ClientConfig{
		User: user.Username,
		Auth: []ssh.AuthMethod{
			// Use a callback rather than PublicKeys so we only consult the
			// agent once the remote server wants it.
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: hostKeyCallback,
	}
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}
	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}
	err = session.Wait()
	if err != nil {
		log.Fatalf("session wait failed: %v", err)
	}
}

func hostKeyCallback(host string, remote net.Addr, key ssh.PublicKey) error {
	// OpenSSH doesn't add :22 to known_hosts file and we want to stay compatible.
	if strings.HasSuffix(host, ":22") {
		host = host[:len(host)-3]
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	rest, err := ioutil.ReadFile(filepath.Join(home, ".ssh", "known_hosts"))
	if err != nil {
		return err
	}
	for {
		var (
			hosts []string
			pk    ssh.PublicKey
		)
		_, hosts, pk, _, rest, err = ssh.ParseKnownHosts(rest)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		for _, h := range hosts {
			if host == h && pk.Type() == key.Type() && bytes.Equal(pk.Marshal(), key.Marshal()) {
				fmt.Printf("%v -> %v\n", hosts, pk)
				return nil
			}
		}
	}
	return fmt.Errorf("host key not found")
}
