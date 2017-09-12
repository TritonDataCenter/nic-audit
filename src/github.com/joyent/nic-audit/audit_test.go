/*
 * Copyright (c) 2017, Joyent, Inc. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
package main

import (
	"github.com/joyent/triton-go/compute"
	"testing"
)

func TestCountOfMatchingNetworkIdsMatchesCIDR(t *testing.T) {
	instanceNetworks := []string{
		"70294144-7680-43d2-9ed0-897ce1658f80",
		"14323a83-b0e3-44e8-bd67-fc7078cc94ba",
		"91ddcc19-b7f9-47b8-8258-f2741bd44112",
	}

	ips := []string{
		"192.168.0.7", "165.122.33.44", "10.2.45.234",
	}

	instance := compute.Instance{
		Networks: instanceNetworks,
		IPs:      ips,
	}

	search := []string{
		"4167e82f-2bd8-46c0-ad4b-7899398c8720",
		"165.122.33.0/21",
	}

	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	count := countMatchingNetworkIds(instance, search, privateBlocks)

	if count != 1 {
		t.Errorf("Expected 1 networks matched. Actually matched %v networks.",
			count)
	}
}

func TestCountOfMatchingNetworkIdsMatchesPublic(t *testing.T) {
	instanceNetworks := []string{
		"70294144-7680-43d2-9ed0-897ce1658f80",
		"14323a83-b0e3-44e8-bd67-fc7078cc94ba",
		"91ddcc19-b7f9-47b8-8258-f2741bd44112",
	}

	ips := []string{
		"192.168.0.7", "165.122.33.44", "10.2.45.234",
	}

	instance := compute.Instance{
		Networks: instanceNetworks,
		IPs:      ips,
	}

	search := []string{
		"4167e82f-2bd8-46c0-ad4b-7899398c8720",
		"public",
	}

	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	count := countMatchingNetworkIds(instance, search, privateBlocks)

	if count != 1 {
		t.Errorf("Expected 1 networks matched. Actually matched %v networks.",
			count)
	}
}

func TestCountOfMatchingNetworkIdsMatchesUUID(t *testing.T) {
	instanceNetworks := []string{
		"70294144-7680-43d2-9ed0-897ce1658f80",
		"14323a83-b0e3-44e8-bd67-fc7078cc94ba",
		"91ddcc19-b7f9-47b8-8258-f2741bd44112",
	}

	ips := []string{
		"192.168.0.7", "165.122.33.44", "10.2.45.234",
	}

	instance := compute.Instance{
		Networks: instanceNetworks,
		IPs:      ips,
	}

	search := []string{
		"14323a83-b0e3-44e8-bd67-fc7078cc94ba",
		"105.111.22.0/21",
	}

	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	count := countMatchingNetworkIds(instance, search, privateBlocks)

	if count != 1 {
		t.Errorf("Expected 1 networks matched. Actually matched %v networks.",
			count)
	}
}

func TestCountOfMatchingNetworkIdsNoMatches(t *testing.T) {
	instanceNetworks := []string{
		"70294144-7680-43d2-9ed0-897ce1658f80",
		"14323a83-b0e3-44e8-bd67-fc7078cc94ba",
		"91ddcc19-b7f9-47b8-8258-f2741bd44112",
	}

	ips := []string{
		"192.168.0.7", "165.122.33.44", "10.2.45.234",
	}

	instance := compute.Instance{
		Networks: instanceNetworks,
		IPs:      ips,
	}

	search := []string{
		"4167e82f-2bd8-46c0-ad4b-7899398c8720",
		"105.111.22.0/21",
	}

	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	count := countMatchingNetworkIds(instance, search, privateBlocks)

	if count != 0 {
		t.Errorf("Expected 0 networks matched. Actually matched %v networks.",
			count)
	}
}

func TestCountOfMatchingNetworkIdsMatchesPublicAndCIDRShouldCountAsOne(t *testing.T) {
	instanceNetworks := []string{
		"70294144-7680-43d2-9ed0-897ce1658f80",
		"14323a83-b0e3-44e8-bd67-fc7078cc94ba",
		"91ddcc19-b7f9-47b8-8258-f2741bd44112",
	}

	ips := []string{
		"192.168.0.7", "165.122.33.44", "10.2.45.234",
	}

	instance := compute.Instance{
		Networks: instanceNetworks,
		IPs:      ips,
	}

	search := []string{
		"public",
		"165.122.33.0/21",
	}

	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	count := countMatchingNetworkIds(instance, search, privateBlocks)

	if count != 1 {
		t.Errorf("Expected 1 networks matched. Actually matched %v networks.",
			count)
	}
}

func TestCountOfMatchingNetworkIdsMatchesTwoCIDRShouldCountAsOne(t *testing.T) {
	instanceNetworks := []string{
		"70294144-7680-43d2-9ed0-897ce1658f80",
		"14323a83-b0e3-44e8-bd67-fc7078cc94ba",
		"91ddcc19-b7f9-47b8-8258-f2741bd44112",
		"14323a83-b0e3-44e8-bd67-fc7078cc94ba",
	}

	ips := []string{
		"192.168.0.7", "165.122.33.44", "10.2.45.234", "165.122.33.22",
	}

	instance := compute.Instance{
		Networks: instanceNetworks,
		IPs:      ips,
	}

	search := []string{
		"165.122.33.0/21",
	}

	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	count := countMatchingNetworkIds(instance, search, privateBlocks)

	if count != 1 {
		t.Errorf("Expected 1 networks matched. Actually matched %v networks.",
			count)
	}
}
