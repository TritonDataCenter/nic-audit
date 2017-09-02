/*
 * Copyright (c) 2017, Joyent, Inc. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
package main

import (
	"context"
	"log"
	"net"
)

import (
	"github.com/joyent/triton-go/compute"
	"github.com/twinj/uuid"
)

// removeNICsBasedOnNetworks removes all NICs from the specified instance
// where the NIC connects to one of the specified networks.
func removeNICsBasedOnNetworks(networks []string,
	instance compute.Instance, client compute.ComputeClient) ([]string, error) {

	listNICsInput := compute.ListNICsInput{
		InstanceID: instance.ID,
	}

	nics, nicsErr := client.Instances().ListNICs(context.Background(), &listNICsInput)

	if nicsErr != nil {
		return nil, nicsErr
	}

	macs := make([]string, len(nics))
	networksToRemove := make([]string, len(nics))
	macCount := 0

	for _, nic := range nics {
		for _, network := range networks {
			if len(network) < 1 {
				continue
			}

			// If our "network" is another UUID it is a simple match
			_, uuidErr := uuid.Parse(network)
			if uuidErr == nil && nic.Network == network {
				macs[macCount] = nic.MAC
				networksToRemove[macCount] = network
				macCount++
				continue
			}

			// If our "network" is a CIDR
			_, ipNet, ipErr := net.ParseCIDR(network)
			if ipErr == nil {
				nicIp := net.ParseIP(nic.IP)
				if ipNet.Contains(nicIp) {
					macs[macCount] = nic.MAC
					networksToRemove[macCount] = network
					macCount++
					continue
				}
			}

			// If our "network" is generalized "public" network
			if network == "public" {
				nicIp := net.ParseIP(nic.IP)
				if isPublicIP(nicIp) {
					macs[macCount] = nic.MAC
					networksToRemove[macCount] = network
					macCount++
					continue
				}
			}
		}
	}

	for i := 0; i < macCount; i++ {
		mac := macs[i]
		network := networksToRemove[i]
		removeNICInput := compute.RemoveNICInput{
			InstanceID: instance.ID,
			MAC:        mac,
		}
		log.Printf("Removing NIC for network [%v] with MAC [%v] from instance [%v]\n",
			mac, network, instance.ID)
		removeErr := client.Instances().RemoveNIC(context.Background(), &removeNICInput)

		if removeErr != nil {
			return nil, removeErr
		}
	}

	return networksToRemove, nil
}
