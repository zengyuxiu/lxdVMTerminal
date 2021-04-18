package main

import "flag"

func main() {
	var host string
	flag.StringVar(&host, "host", "192.168.31.14", "host ip address")
	flag.Parse()
	StartServer(host)
}
