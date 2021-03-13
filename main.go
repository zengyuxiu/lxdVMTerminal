package main

import (
	"fmt"
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/logger"
	"io/ioutil"
)

func main()  {
	d,err  := InitLxdInstanceServer("192.168.31.14")
	if err != nil{
		logger.Error(err.Error())
		return
	}
	err = vga(*d,"win7")
	if err != nil{
		logger.Error(err.Error())
		return
	}
}

func InitLxdInstanceServer(manageip string) (*lxd.InstanceServer, error) {
	cert, err := ioutil.ReadFile("/home/zyx/.config/lxc/client.crt")
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile("/home/zyx/.config/lxc/client.key")
	if err != nil {
		return nil, err
	}

	args := &lxd.ConnectionArgs{
		TLSClientCert:      string(cert),
		TLSClientKey:       string(key),
		InsecureSkipVerify: true,
	}
	server, err := lxd.ConnectLXD(fmt.Sprintf("https://%s:8443", manageip), args)
	if err != nil {
		return nil, err
	}

	return &server, nil
}