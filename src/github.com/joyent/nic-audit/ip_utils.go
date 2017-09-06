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
