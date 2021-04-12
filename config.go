package main

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	} `yaml:"server"`
}
