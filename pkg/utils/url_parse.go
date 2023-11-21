package utils

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
)

func ParseURIToHostPort(uri string) (ip string, port int, err error) {
	rawUrl, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", 0, fmt.Errorf("cant parse URI: %s to url struct, %w", uri, err)
	}

	if rawUrl.Port() == "" {
		rawUrl.Host = rawUrl.Host + ":80"
	}

	host := rawUrl.Host
	ip, portString, err := net.SplitHostPort(host)
	if err != nil {
		return "", 0, fmt.Errorf("cant parse URI: %s, %w", uri, err)
	}

	if net.ParseIP(ip) == nil {
		return "", 0, fmt.Errorf("cant parse URI: %s, invalid IP address format", uri)
	}

	port, err = strconv.Atoi(portString)
	if err != nil {
		return "", 0, fmt.Errorf("cant parse URI: %s to url struct, %w", uri, err)
	}

	return ip, port, nil
}
