package config

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	cluster_config := "../clusters/sample.conf"
	config := ReadConfig(&cluster_config)

	if config.Storage.DataPath != "/data" {
		t.Errorf("%s != %s", config.Storage, "/data")
	}
	if config.Credentials.RestUsername != "Administrator" {
		t.Errorf("%s != %s", config.Credentials.RestUsername, "Administrator")
	}
	if config.Servers.CB[0] != "192.168.1.1" {
		t.Errorf("%s != %s", config.Servers.CB[0], "192.168.1.1")
	}
}
