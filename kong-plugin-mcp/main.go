package main

import (
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

const Version = "0.1"
const Priority = 1000

/*
This plugin forwards requests to an MCP server.

Currently, it is unnecessary, as Kong natively supports HTTP proxies. However, a goal of this project is to learn more about Kong plugins. See the README for more information.
*/
func main() {
	server.StartServer(New, Version, Priority)
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	handleHTTPProxy(kong, conf)
}
