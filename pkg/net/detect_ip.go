package net

import (
	"errors"
	"fmt"
	"net"
)

var (
	ErrNotDetected = errors.New("couldn't detect ip")
	ErrAFewIPs     = errors.New("a few ip were detected")
)

func detectIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("InterfaceAddrs: %w", err)
	}

	var suitableIPs []string

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				suitableIPs = append(suitableIPs, ipNet.IP.String())
			}
		}
	}
	switch len(suitableIPs) {
	case 0:
		return "", ErrNotDetected
	case 1:
		return suitableIPs[0], nil
	default:
		return "", ErrAFewIPs
	}
}
