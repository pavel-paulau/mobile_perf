package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/pavel-paulau/mobile_perf/config"
	"github.com/pavel-paulau/mobile_perf/helpers"
)

type Build struct {
	Version, Arch, Pkg string
}

type SyncGatewayInstaller struct {
	build          Build
	cluster_config config.ClusterConfig
}

func (i *SyncGatewayInstaller) generatePackageMeta() (fname, url string) {
	fname = fmt.Sprintf("couchbase-sync-gateway-community_%s_%s.%s",
		i.build.Version, i.build.Arch, i.build.Pkg)
	url = fmt.Sprintf("http://packages.couchbase.com/releases/couchbase-sync-gateway/%s/%s",
		i.build.Version, fname)
	resp, err := http.Head(url)
	if err != nil || resp.Status != "200 OK" {
		log.Fatalf("Target Sync Gateway build not found")
	}
	log.Printf("Found %s", url)
	return
}

func GetFlags() (config_file, version string) {
	flag.StringVar(&config_file, "cluster", "", "Path to cluster configuration")
	flag.StringVar(&version, "version", "", "Sync Gateway version")
	flag.Parse()

	if flag.NFlag() != 2 {
		flag.PrintDefaults()
		log.Fatal("Missing arguments")
	}
	return
}

func main() {
	remote := helpers.RemoteHelper{}
	remote.Init("root", "couchbase")

	config_file, version := GetFlags()

	cluster_config := config.ReadConfig(&config_file)

	arch := remote.DetectArch(cluster_config.Servers.SG[0])
	pkg := remote.DetectPkg(cluster_config.Servers.SG[0])
	build := Build{version, arch, pkg}

	installer := SyncGatewayInstaller{build, cluster_config}

	fname, url := installer.generatePackageMeta()

	for _, host := range cluster_config.Servers.SG {
		remote.UninstallPackage(host, fname, pkg)
	}

	for _, host := range cluster_config.Servers.SG {
		remote.Wget(host, url, "/tmp")
	}

	for _, host := range cluster_config.Servers.SG {
		remote.InstallPackage(host, fname, pkg)
	}
}
