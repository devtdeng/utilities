package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// xrocss compile
// GOOS=linux GOARCH=amd64 go build -o golang-dns-<VERION> main.go
// GOOS=darwin GOARCH=amd64 go build -o golang-dns-<VERION> main.go

// export GODEBUG=netdns=2
// export GODEBUG=netdns=cgo+2
// export GODEBUG=netdns=go+2
// ./golang-dns-<VERION> <HOSTNAME>

func main() {
	hostName := os.Args[1]
	// dnsServer := os.Args[2]

	t := time.Now()
	fmt.Println("DNS resolution starts at: " + t.Format("20060102150405"))
	ips, err := net.LookupIP(hostName)
	t = time.Now()
	fmt.Println("DNS resolution finishes at: " + t.Format("20060102150405"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	for _, ip := range ips {
		fmt.Printf("%s. IN A %s\n", hostName, ip.String())
	}
}
