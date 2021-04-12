package main

import (
	"fmt"
	lxd "github.com/lxc/lxd/client"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func InitLxdInstanceServer(host string) (*lxd.InstanceServer, error) {
	ConfigFile, err := os.Open("config.yml")
	if err != nil {
		return nil, err
	}
	defer ConfigFile.Close()

	var cfg Config
	decoder := yaml.NewDecoder(ConfigFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	cert, err := ioutil.ReadFile(cfg.Server.Cert)
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile(cfg.Server.Key)
	if err != nil {
		return nil, err
	}

	args := &lxd.ConnectionArgs{
		TLSClientCert:      string(cert),
		TLSClientKey:       string(key),
		InsecureSkipVerify: true,
	}
	server, err := lxd.ConnectLXD(fmt.Sprintf("https://%s:%s", host, cfg.Server.Port), args)
	if err != nil {
		return nil, err
	}

	return &server, nil
}
