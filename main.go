package main

import (
	"feather/config"
	"feather/network"
	"feather/repository"
	"feather/service"
	"flag"
)

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", "localhost:8080", "port set")

func main() {
	flag.Parse()

	c := config.NewConfig(*pathFlag)
	if repository, err := repository.NewRepository(c); err != nil {
		panic(err)
	} else {
		n := network.NewServer(service.NewService(repository), *port)
		if err := n.StartServer(); err != nil {
			panic(err)
		}
	}
}
