package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/asaskevich/govalidator"
	"github.com/glendc/go-external-ip"
	"net"
	"os"
)

func parse_args() string {
	parser := argparse.NewParser("vpn-twitcher", "checks that a given ip is not the same as the public ip")
	actual_public_ip := parser.String("i", "ip_address", &argparse.Options{Required: true, Help: "public ip/domain to check"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	return *actual_public_ip
}
func process_ip(actual_public_ip string) (string, int) {
	/*
	 * returns (ipv4, return_code)
	 */
	if !govalidator.IsIPv4(actual_public_ip) {
		// this might be a domain name
		ips, lookup_err := net.LookupIP(actual_public_ip)
		if lookup_err != nil {
			fmt.Println("Failed getting IPv4 from " + actual_public_ip)
			return "", 2
		}
		for _, ip := range ips {
			if govalidator.IsIPv4(ip.String()) {
				actual_public_ip = ip.String()
			}
		}
		if !govalidator.IsIPv4(actual_public_ip) {
			fmt.Println("Failed getting IPv4 from " + actual_public_ip)
			return "", 2
		}
	} else if govalidator.IsIPv6(actual_public_ip) {
		fmt.Println("IPv6 is currently not supported")
		return "", 1
	}
	return actual_public_ip, 0
}
func main() {
	actual_public_ip := parse_args()

	ipv4, return_code := process_ip(actual_public_ip)
	if return_code != 0 {
		os.Exit(return_code)
	}

	// get public IP
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		fmt.Println("failed getting public ip")
		os.Exit(2)
	}
	if ip.String() == ipv4 {
		// DO SOMETHING VPN IS NOT WORKING CORRECTLY!
		os.Exit(3)
	}
}
