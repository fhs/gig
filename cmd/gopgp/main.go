package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/openpgp"
)

// delta=; git verify-commit HEAD
// /usr/bin/gpg --keyid-format=long --status-fd=1 --verify /tmp/.git_vtag_tmpPTg5sV -
//
// github pgp key:  https://github.com/fhs.gpg

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	ring, err := openpgp.ReadKeyRing(f)
	if err != nil {
		panic(err)
	}
	for _, ent := range ring {
		fmt.Printf("ent = %v\n", ent)
	}
}

func write() {
	ent, err := openpgp.NewEntity("Fazlul Shahriar", "git", "fshahriar@gmail.com", nil)
	if err != nil {
		panic(err)
	}
	err = ent.SerializePrivate(os.Stdout, nil)
	if err != nil {
		panic(err)
	}
}
