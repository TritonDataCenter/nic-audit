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
	"net"
	"testing"
)

func TestParseMultipleCIDRsCanParseASingleValidCIDR(t *testing.T) {
	cidr := "172.16.0.0/12"
	_, expected, _ := net.ParseCIDR(cidr)

	cidrs, err := parseMultipleCIDRs(cidr)

	if err != nil {
		t.Error(err)
	}

	if len(cidrs) != 1 {
		t.Errorf("Too many elements returned. Expected 1. Actually: %v",
			len(cidrs))
	}

	if cidrs[0].String() != expected.String() {
		t.Errorf("Unexpected value: %v", cidrs[0])
	}
}

func TestParseMultipleCIDRsWillErrorOnASingleInvalidCIDR(t *testing.T) {
	cidr := "172.16.0.0/99"
	_, err := parseMultipleCIDRs(cidr)

	if err == nil {
		t.Error("Expected error and none was thrown")
	}
}

func TestParseMultipleCIDRsCanParseMultipleValidCIDRs(t *testing.T) {
	input := "172.16.0.0/12,192.168.24.0/21,10.0.0.0/8"
	expected := "[{172.16.0.0 fff00000} {192.168.24.0 fffff800} {10.0.0.0 ff000000}]"
	cidrs, err := parseMultipleCIDRs(input)

	if err != nil {
		t.Error(err)
	}

	if len(cidrs) != 3 {
		t.Errorf("Incorrect number of elements returned. Expected 3. Actually: %v",
			len(cidrs))
	}

	cidrsAsString := fmt.Sprintf("%v", cidrs)

	if cidrsAsString != expected {
		t.Errorf("Unexpected value: %v\nExpected: %v", cidrsAsString, expected)
	}
}

func TestParseMultipleCIDRsCanParseMultipleValidCIDRsWithWhitespace(t *testing.T) {
	input := " 172.16.0.0/12,    192.168.24.0/21,   10.0.0.0/8"
	expected := "[{172.16.0.0 fff00000} {192.168.24.0 fffff800} {10.0.0.0 ff000000}]"
	cidrs, err := parseMultipleCIDRs(input)

	if err != nil {
		t.Error(err)
	}

	if len(cidrs) != 3 {
		t.Errorf("Incorrect number of elements returned. Expected 3. Actually: %v",
			len(cidrs))
	}

	cidrsAsString := fmt.Sprintf("%v", cidrs)

	if cidrsAsString != expected {
		t.Errorf("Unexpected value: %v\nExpected: %v", cidrsAsString, expected)
	}
}

func TestParseMultipleCIDRsCanParseMultipleValidAndOneInvalidCIDRs(t *testing.T) {
	input := " 172.16.0.0/12,192.168.24.0/99,10.0.0.0/8"
	expected := "[{172.16.0.0 fff00000} {10.0.0.0 ff000000}]"
	cidrs, err := parseMultipleCIDRs(input)

	if err != nil {
		t.Error(err)
	}

	if len(cidrs) != 2 {
		t.Errorf("Incorrect number of elements returned. Expected 2. Actually: %v",
			len(cidrs))
	}

	cidrsAsString := fmt.Sprintf("%v", cidrs)

	if cidrsAsString != expected {
		t.Errorf("Unexpected value: %v\nExpected: %v", cidrsAsString, expected)
	}
}

func TestParseMultipleCIDRsCanParseMultipleValidCIDRsWithDuplicate(t *testing.T) {
	input := "172.16.0.0/12,192.168.24.0/21,172.16.0.0/12"
	expected := "[{172.16.0.0 fff00000} {192.168.24.0 fffff800}]"
	cidrs, err := parseMultipleCIDRs(input)

	if err != nil {
		t.Error(err)
	}

	if len(cidrs) != 2 {
		t.Errorf("Incorrect number of elements returned. Expected 2. Actually: %v",
			len(cidrs))
	}

	cidrsAsString := fmt.Sprintf("%v", cidrs)

	if cidrsAsString != expected {
		t.Errorf("Unexpected value: %v\nExpected: %v", cidrsAsString, expected)
	}
}
