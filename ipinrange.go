package main

// ipinrange :
//
// Basic ipV4 extrator base on network range
//
// checkout similar code  in extractip project
//
// How it works :
//  read stdin, output filteredinput if error to output
//
// Evolutions :
// 2023/07/20 : V0.1
// 2023/11/24 : Negative option

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
)

func main() {

	var (
		negativeFlag bool
		subnet       *net.IPNet
		E            error
	)

	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	// Same things for others
	// emails := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	// DomainUrl := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)
	// words := regexp.MustCompile(`[\p{L}]+`) // Without numbers
	// words := regexp.MustCompile("\\P{M}+") // With numbers ?

	switch {
	case len(os.Args) == 3 && os.Args[1] == "-N":
		negativeFlag = true
		_, subnet, E = net.ParseCIDR(os.Args[2])
	case len(os.Args) == 2:
		_, subnet, E = net.ParseCIDR(os.Args[1])
	default:
		fmt.Printf("Usage: ipinrange [-N(egative)] network/x")
		os.Exit(-1)
	}

	if E != nil {
		fmt.Fprint(os.Stderr, E)
		os.Exit(-1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()

		submatchall := re.FindAllString(text, -1)
		for _, element := range submatchall {
			contains := subnet.Contains(net.ParseIP(element))
			if (contains && !negativeFlag) || (negativeFlag && !contains) {
				fmt.Println(text)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}

// network := "192.168.5.0/24"
// clientips := []string{
//     "192.168.5.1",
//     "192.168.6.0",
// }
// _, subnet, _ := net.ParseCIDR(network)
// for _, clientip := range clientips {
//     ip := net.ParseIP(clientip)
//     if subnet.Contains(ip) {
//         fmt.Println("IP in subnet", clientip)
//     }
// }
