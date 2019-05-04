package util

import (
	"net"
)

type hostport = string

func MustGetIp(hp hostport) (ip, port string) {
	if hp == "" {
		panic("received emptys string as hostport")
	}
	ip, port, err := net.SplitHostPort(hp)
	if err != nil {
		panic(err)
	}
	return ip, port
}
