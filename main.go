package main

import (
	"flag"
	"fmt"
)
import (
	"net_transfer/client"
	"net_transfer/server"
)

func main() {
	op := flag.String("op", "server", "")
	ip := flag.String("ip", "0.0.0.0", "")
	port := flag.String("port", "8888", "")
	flag.Parse()
	//server.StartProxy(*ip, *port)
	//server.StartServer(*ip, *port)
	//client.StartClient(*ip, *port)
	if *op == "server" {
		server.StartServer(*ip, *port)
	} else if *op == "client" {
		fmt.Println("start client")
		client.StartClient(*ip, *port)
	} else if *op == "proxy" {
		server.StartProxy(*ip, *port)
	}
}
