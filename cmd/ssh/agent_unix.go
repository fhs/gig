// +build !plan9

package main

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh/agent"
)

func newAgentClient() (agent.ExtendedAgent, error) {
	// ssh-agent(1) provides a UNIX socket at $SSH_AUTH_SOCK.
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return nil, fmt.Errorf("open $SSH_AUTH_SOCK: %v", err)
	}
	return agent.NewClient(conn), nil
}
