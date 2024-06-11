package util

import (
	"net"
)

func GetIPv4Address() string {
	// get all network interfaces
	interfaces, err := net.Interfaces()
	// log.Println("DEBUG | interfaces:", interfaces)
	if err != nil {
		return ""
	}

	// iterate over the interfaces
	for _, iface := range interfaces {
		// get the addresses associated with the interface
		addresses, err := iface.Addrs()
		// log.Println("DEBUG | addresses:", addresses)
		if err != nil {
			continue
		}

		// iterate over the addresses
		for _, address := range addresses {
			// log.Println("DEBUG | address:", address)
			// check if the address is an IPv4 address
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}

	return ""
}
