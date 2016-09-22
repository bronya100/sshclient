package sshclient

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

var (
	host, username, password, keyfile string
	keyauth                           ssh.AuthMethod
)

func init() {
	host = os.Getenv("SSH_HOST")
	if len(host) == 0 {
		panic("Error: SSH_HOST not set")
	}
	host += ":22"

	username = os.Getenv("SSH_USERNAME")
	if len(username) == 0 {
		panic("Error: SSH_USERNAME not set")
	}

	password = os.Getenv("SSH_PASSWORD")
	if len(password) == 0 {
		panic("Error: SSH_PASSWORD not set")
	}

	keyfile = os.Getenv("SSH_PRIVATE")
	if len(keyfile) == 0 {
		panic("Error: SSH_PRIVATE not set (private keyfile)")
	}
}

func TestSSHKey(t *testing.T) {
	var err error
	keyauth, err = KeyAuth(keyfile)
	if err != nil {
		t.Fatal("keyauth error:", err)
	}
	client, err := DialSSH(host, username, 5, keyauth)
	if err != nil {
		t.Fatal("keyauth dial error:", err)
	}
	cmd := "uptime"
	r := Run(client, cmd)
	if r.Err != nil {
		t.Fatal("keyauth run error:", err)
	}
}

func TestSSHKeyauth(t *testing.T) {
	client, err := DialKey(host, username, keyfile, 5)
	if err != nil {
		t.Fatal("keyauth dial error:", err)
	}
	cmd := "logname"
	r := Run(client, cmd)
	if r.Err != nil {
		t.Fatal("keyauth run error:", err)
	}
	if strings.TrimSpace(r.Stdout) != username {
		t.Fatal("keyauth command failed. expected", username, "got", r.Stdout)
	}
}

func TestSSHClient(t *testing.T) {
	cmd := "hostname"
	timeout := 5
	rc, stdout, stderr, err := Exec(host, username, password, cmd, timeout)
	if err != nil {
		t.Error("ssh connect error:", err)
	}
	if rc > 0 {
		t.Error("ssh execution error:", stderr)
	} else if len(stderr) > 0 {
		t.Error("ssh execution error:", stderr)
	} else {
		t.Log("client returned:", stdout)
	}
}

func TestSSHStderr(t *testing.T) {
	cmd := "lsX"
	timeout := 5
	_, stdout, stderr, _ := Exec(host, username, password, cmd, timeout)
	if len(stdout) > 0 {
		t.Log("ssh stdout", stdout)
	}
	if len(stderr) > 0 {
		t.Log("ssh stderr:", stderr)
	}
}

func TestSSHTimeout(t *testing.T) {
	cmd := "sleep 10"
	timeout := 5
	rc, _, stderr, err := Exec(host, username, password, cmd, timeout)
	if err == nil {
		t.Error("ssh timeout failed")
	}
	if rc > 0 {
		t.Error("ssh execution error:", stderr)
	} else if len(stderr) > 0 {
		t.Error("ssh execution error:", stderr)
	}
}
