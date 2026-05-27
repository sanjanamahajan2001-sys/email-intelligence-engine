package main

import (
	"fmt"
	"net"
)

func main() {
	mxs, err := net.LookupMX("mailinator.com")
	if err != nil {
		fmt.Printf("Error looking up MX: %v\n", err)
		return
	}
	for _, mx := range mxs {
		fmt.Printf("MX: %s\n", mx.Host)
	}
}
