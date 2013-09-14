package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Servers struct {
	CB, SG []string
	LB     string
}

type Storage struct {
	DataPath, IndexPath string
}

type Credentials struct {
	RestUsername, RestPassword, SshUsername, SshPassword string
}

type ClusterConfig struct {
	Servers     Servers
	Storage     Storage
	Credentials Credentials
}

func ReadConfig(config_file *string) (config ClusterConfig) {
	cluster_config, err := ioutil.ReadFile(*config_file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(cluster_config, &config)
	if err != nil {
		log.Fatal(err, string(cluster_config))
	}
	return
}
