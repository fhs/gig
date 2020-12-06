// Copyright 2020 Fazlul Shahriar. All rights reserved.
// Use of this source code is governed by the
// license that can be found in the LICENSE file.

package cli

import (
	"crypto/rand"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"unsafe"
)

func setupFactotum(t *testing.T) {
	if _, err := rfork(syscall.RFNAMEG); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("/bin/auth/factotum", "-n")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	// These keys were generated on Plan 9 using: auth/rsagen -t 'service=ssh'
	keys, err := ioutil.ReadFile("testdata/factotum.keys")
	if err != nil {
		t.Fatal(err)
	}
	w, err := os.OpenFile("/mnt/factotum/ctl", os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	for _, line := range strings.Split(string(keys), "\n") {
		if _, err := w.WriteString(line); err != nil {
			t.Fatal(err)
		}
	}
}

func TestFactotum(t *testing.T) {
	setupFactotum(t)
	signers, err := factotumSigners()
	if err != nil {
		t.Fatal(err)
	}
	for i, priv := range signers {
		pub := priv.PublicKey()
		t.Logf("checking key %v: %v\n", i, priv)
		data := []byte("sign me")
		sig, err := priv.Sign(rand.Reader, data)
		if err != nil {
			t.Fatalf("Sign(%T): %v", priv, err)
		}
		if err := pub.Verify(data, sig); err != nil {
			t.Errorf("publicKey.Verify(%T): %v", priv, err)
		}
		sig.Blob[5]++
		if err := pub.Verify(data, sig); err == nil {
			t.Errorf("publicKey.Verify on broken sig did not fail")
		}
	}
}

func rfork(flags int) (pid int, err error) {
	r1, _, _ := syscall.RawSyscall(syscall.SYS_RFORK, uintptr(flags), 0, 0)
	if int32(r1) == -1 {
		return 0, syscall.NewError(errstr())
	}
	return int(r1), nil
}

func cstring(s []byte) string {
	for i, b := range s {
		if b == 0 {
			return string(s[0:i])
		}
	}
	return string(s)
}

func errstr() string {
	var buf [syscall.ERRMAX]byte

	syscall.RawSyscall(syscall.SYS_ERRSTR, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), 0)

	buf[len(buf)-1] = 0
	return cstring(buf[:])
}
