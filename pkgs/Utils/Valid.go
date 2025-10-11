package Utils

import (
	"net"
	"strconv"
	"strings"
)

func IsValidIP(ipString string) bool {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return false
	}
	if ip.To4() != nil {
		return true // IPv4
	}
	if ip.To16() != nil {
		return true // IPv6
	}
	return false
}

func IsHost(host string) bool {
	if IsValidIP(host) {
		return true
	}
	_, err := net.LookupHost(host)
	if err != nil {
		return false
	}
	return true
}

func IsValidPort(portString string) bool {
	portString = strings.TrimSpace(portString)
	if portString == "" || len(portString) > 5 {
		return false
	}
	for _, c := range portString {
		if c < '0' || c > '9' {
			return false
		}
	}
	if strings.HasPrefix(portString, "0") && len(portString) > 1 {
		return false
	}

	portNum, err := strconv.Atoi(portString)
	if err != nil || portNum < 1 || portNum > 65535 {
		return false
	}
	return true
}

func IsValidBrokers(addr, offset string) bool {
	if addr == "" {
		return false
	}
	if offset != "" {
		offset = ","
	}
	adders := strings.Split(addr, offset)
	for _, v := range adders {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if strings.Contains(v, ":") {
			hostPort := strings.SplitN(v, ":", 2)
			if len(hostPort) != 2 || !IsHost(hostPort[0]) || !IsValidPort(hostPort[1]) {
				return false
			}
		} else {
			if !IsHost(v) {
				return false
			}
			if offset == "" || !IsValidPort(offset) {
				return false
			}
		}
	}
	return true
}
