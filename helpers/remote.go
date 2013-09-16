package helpers

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"code.google.com/p/go.crypto/ssh"

	"github.com/pavel-paulau/mobile_perf/config"
)

type clientPassword string

func (p clientPassword) Password(user string) (string, error) {
    return string(p), nil
}

type RemoteHelper struct {
	cluster_config *config.ClusterConfig
	config *ssh.ClientConfig
}

func (h *RemoteHelper) Init(username, password string) {
	h.config = &ssh.ClientConfig{
        User: username,
        Auth: []ssh.ClientAuth{
            ssh.ClientAuthPassword(clientPassword(password)),
        },
    }
} 

func (h *RemoteHelper) runCmd(host, cmd string) string {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), h.config)
    if err != nil {
        log.Fatalf("Failed to dial: " + err.Error())
    }
    defer client.Close()

    session, err := client.NewSession()
    if err != nil {
        log.Fatal("unable to create session: %s", err)
    }
    defer session.Close()

    err = session.RequestPty("xterm", 80, 40, ssh.TerminalModes{});
    if err != nil {
        log.Fatal("request for pseudo terminal failed: %s", err)
    }

    var b bytes.Buffer
    session.Stdout = &b

    err = session.Run(cmd)
    if err != nil {
        log.Fatal("Failed to run: " + err.Error())
    }

    return strings.Replace(strings.Replace(b.String(), "\n", "", -1), "\r", "", -1)
}

func (h *RemoteHelper) DetectArch(host string) string {
	return h.runCmd(host, "uname -i")
}

func (h *RemoteHelper) DetectPkg(host string) string {
	os := h.runCmd(host, "python -c 'import platform; print platform.dist()[0]'")
    if os == "Ubuntu" {
    	return "deb"
    }
    return "rpm"
}

func (h *RemoteHelper) Wget(host, url, outdir string) {
	log.Printf("Fetching %s", url)

	cmd := fmt.Sprintf("wget -nc '%s' -P %s", url, outdir)
	h.runCmd(host, cmd)
}

func (h *RemoteHelper) UninstallPackage(host, fname, pkg string) {
	log.Printf("Uninstalling Sync Gateway")

	var cmd string
	if pkg == "deb" {
		cmd = "yes | sudo apt-get remove couchbase-sync-gateway"
    } else {
        cmd = "yes | sudo yum remove couchbase-sync-gateway"
    }
	h.runCmd(host, cmd)
}

func (h *RemoteHelper) InstallPackage(host, fname, pkg string) {
	log.Printf("Installing Sync Gateway")

	var cmd string
	if pkg == "deb" {
		cmd = fmt.Sprintf("yes | sudo dpkg -i /tmp/%s", fname)
    } else {
        cmd = fmt.Sprintf("yes | rpm -i /tmp/%s", fname)
    }
	h.runCmd(host, cmd)
}
