package main

import (
	"flag"
	"runtime"
)

import (
	"httpflash/common"
	"httpflash/httptask"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&common.URL, "url", "127.0.0.1:80", "<ipaddress:port> | default 127.0.0.1:80")
	flag.StringVar(&common.METHOD, "method", "GET", "<GET | POST | PUT | HEAD | PATCH | DELETE | OPTIONS | TRACE | CONNECT> | default GET")
	flag.StringVar(&common.HEADER, "header", "", "<'key: value, key: value'> | default none")
	flag.StringVar(&common.COOKIE, "cookie", "", "<'name=value;name=value'> | default none")
	flag.StringVar(&common.DATA, "data", "", "<'name=value;name=value'> | default none")

	flag.Uint64Var(&common.VCLIENT, "client", 10, "default 10")
	flag.Uint64Var(&common.VCLICK, "click", 10, "default 10")

	flag.Parse()
}

func main() {
	httptask.HttpTask()
}
