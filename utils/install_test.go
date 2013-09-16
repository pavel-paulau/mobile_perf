package main

import (
	"fmt"
	"testing"

	"github.com/pavel-paulau/mobile_perf/config"
)

func TestPackageMeta(t *testing.T) {
	build := Build{"1.0-beta", "x86_64", "deb"}
	cluster_config := config.ClusterConfig{}
	installer := SyncGatewayInstaller{build, cluster_config}

	fname, url := installer.generatePackageMeta()
	expected_fname := "couchbase-sync-gateway-community_1.0-beta_x86_64.deb"
	expected_url := fmt.Sprintf("http://packages.couchbase.com/releases/couchbase-sync-gateway/1.0-beta/%s",
		expected_fname)

	if expected_fname != fname {
		t.Errorf("%s != %s", expected_fname, fname)
	}
	if expected_url != url {
		t.Errorf("%s != %s", expected_url, url)
	}
}
