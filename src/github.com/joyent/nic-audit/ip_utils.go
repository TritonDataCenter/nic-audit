/*
 * Copyright (c) 2017, Joyent, Inc. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// isPrivateIP determines if the specified IP address is on a
// private network.
func isPrivateIP(ip net.IP, privateBlocks []string) bool {
	for i := 0; i < len(privateBlocks); i++ {
		_, privateBlock, parseErr := net.ParseCIDR(privateBlocks[i])

		if parseErr != nil {
			msg := fmt.Sprintf("Unable to parse private network [%v]. %v",
				privateBlock, parseErr)
			log.Fatalln(msg)
		}

		if privateBlock.Contains(ip) {
			return true
		}
	}

	return false
}

// isPublicIP determines if the specified IP address is not a RFC 1918
// private network.
func isPublicIP(ip net.IP, privateBlocks []string) bool {
	return !isPrivateIP(ip, privateBlocks)
}

// parseMultipleCIDRs parses a comma delimited list of CIDRs and returns an
// array containing all valid values.
func parseMultipleCIDRs(input string) ([]net.IPNet, error) {
	if strings.Contains(input, ",") {
		uniqueCidrs := make(map[string]bool)

		elements := strings.Split(input, ",")

		cidrs := make([]net.IPNet, len(elements))

		count := 0
		for _, element := range elements {
			cidr := strings.TrimSpace(element)

			if len(cidr) < 1 {
				continue
			}

			_, ipNet, ipErr := net.ParseCIDR(strings.TrimSpace(cidr))

			if ipErr != nil {
				log.Printf("Invalid CIDR specified: %v\n", cidr)
				continue
			}

			if uniqueCidrs[ipNet.String()] {
				log.Printf("Duplicate CIDR specified: %v\n", cidr)
				continue
			}

			cidrs[count] = *ipNet
			uniqueCidrs[ipNet.String()] = true
			count++
		}

		return cidrs[0:count], nil
	}

	_, ipNet, ipErr := net.ParseCIDR(strings.TrimSpace(input))
	if ipErr != nil {
		return nil, ipErr
	} else {
		return []net.IPNet{*ipNet}, nil
	}
}
