package main

// ipinrange :
//
// Basic ipV4 filter based on network range
//
// checkout similar code  in extractip project
//
// How it works :
//  read stdin, output filtered input
//    if error output to stderr
//
// Evolutions :
// 2023/07/20 : V0.1
// 2023/11/24 : Negative option
// 2023/11/29 : v1.3
//              - remove reserved ip addresses
//              - allow array as arg
//              - know local network
// so you can extract text  :
//  cat logs | ipinrange -N local
//
//

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {

	var (
		negativeFlag bool
		subnetA      []*net.IPNet
		netARg       string
		lenArg       = len(os.Args)
		re           = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

		// https://en.wikipedia.org/wiki/Reserved_IP_addresses
		Reserved_IP_addresses = []string{
			"0.0.0.0/8",
			"10.0.0.0/8",
			"100.64.0.0/10",
			"127.0.0.0/8",
			"169.254.0.0/16",
			"172.16.0.0/12",
			"192.0.0.0/24",
			"192.0.2.0/24",
			"192.88.99.0/24",
			"192.168.0.0/16",
			"198.18.0.0/15",
			"198.51.100.0/24",
			"203.0.113.0/24",
			"224.0.0.0/4",
			"233.252.0.0/24",
			"240.0.0.0/4",
			"255.255.255.255/32",
		}
	)

	switch {
	case lenArg >= 2 && os.Args[1] == "-n":
		negativeFlag = true
		fallthrough

	case lenArg > 1:
		netARg = os.Args[lenArg-1]

	default:
		fmt.Printf("Usage: ipinrange [-n(egative)] []network/x")
		os.Exit(-1)
	}

	if netARg == "local" { // Special Case
		subnetA = parseNetStringtoCIDR(local)
	} else {
		subnetA = argstoCIDR(netARg)
	}

	Reserved_IP_addressesCIDR := parseNetStringtoCIDR(Reserved_IP_addresses)

	scanner := bufio.NewScanner(os.Stdin) // Reading stdin
	for scanner.Scan() {
		text := scanner.Text()
		submatchall := re.FindAllString(text, -1) // Finding ipv4
		for _, element := range submatchall {
			elementIP := net.ParseIP(element)                         // Real ipv4 ?
			found := isIn(elementIP, subnetA)                         // march  our network arg ?
			if (found && !negativeFlag) || (negativeFlag && !found) { // print ?
				if !isIn(elementIP, Reserved_IP_addressesCIDR) { // exclude bogon
					fmt.Println(text)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

// argstoCIDR return CIDRparsed array from "net,net"
func argstoCIDR(arg string) []*net.IPNet {
	s := strings.Split(arg, ",")
	return parseNetStringtoCIDR(s)
}

// isIn check if IP is in a []blocknet
func isIn(ip net.IP, reserved []*net.IPNet) bool {
	found := false
	for _, v := range reserved {
		if v != nil && v.Contains(ip) {
			return true
		}
	}
	return found
}

// parseNetStringtoCIDR convert []string to []*net.IPNet
func parseNetStringtoCIDR(block []string) []*net.IPNet {
	netBlock := make([]*net.IPNet, len(block))
	for k, v := range block {
		_, net, err := net.ParseCIDR(v)
		if err == nil {
			netBlock[k] = net
		}
	}
	return netBlock
}
