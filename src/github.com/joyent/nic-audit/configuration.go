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
	"io"
	"log"
	"net"
	"os"
)

import (
	"github.com/flynn/json5"
	"github.com/twinj/uuid"
)

// Configuration contains all of the configuration values for the application.
type Configuration struct {
	EmailAlerts          EmailAlerts         `json:"email_alerts"`
	PrivateNetworkBlocks []string            `json:"private_network_blocks"`
	NicGroups            map[string][]string `json:"nic_groups"`
	Accounts             []Account           `json:"accounts"`
}

// EmailAlerts contains the configuration needed to send an email to alert when
// an offending network match is found.
type EmailAlerts struct {
	SmtpServer string `json:"smtp_server"`
	To         []string
	CC         []string
	BCC        []string
	From       string
	FromName   string `json:"from_name"`
	Subject    string
}

// Account contains the configuration details describing a single Triton
// account.
type Account struct {
	Description      string
	TritonUrl        string   `json:"triton_url"`
	AccountName      string   `json:"account_name"`
	KeyPath          string   `json:"key_path"`
	KeyId            string   `json:"key_id"`
	NetworksToRemove []string `json:"networks_to_remove"`
}

// readConfigFromFile parses a json5 configuration from the specified path.
func readConfigFromFile(configFile string) (Configuration, error) {
	if !exists(configFile) {
		log.Fatalf("Configuration file [%v] doesn't exist", configFile)
	}

	if !isReadable(configFile) {
		log.Fatalf("Configuration file [%v] is not accessible", configFile)
	}

	reader, fileOpenErr := os.Open(configFile)

	if fileOpenErr != nil {
		return Configuration{}, fileOpenErr
	}

	defer reader.Close()

	return readConfig(reader)
}

// readConfig parses a json5 configuration from the specified reader object.
func readConfig(reader io.Reader) (Configuration, error) {
	decoder := json5.NewDecoder(reader)
	var config Configuration
	decodeErr := decoder.Decode(&config)

	if decodeErr != nil {
		return Configuration{}, decodeErr
	}

	return config, nil
}

// validateConfiguration verifies if a given configuration instance has the
// correct settings.
func validateConfiguration(config Configuration) {
	for _, account := range config.Accounts {
		if !exists(account.KeyPath) {
			msg := fmt.Sprintf("Unable to audit account [%v] because "+
				"private key doesn't exist [%v]", account.AccountName, account.KeyPath)
			log.Fatal(msg)
		}

		if !isReadable(account.KeyPath) {
			msg := fmt.Sprintf("Unable to audit account [%v] because "+
				"private key isn't accessible [%v]", account.AccountName, account.KeyPath)
			log.Fatal(msg)
		}

		for _, network := range account.NetworksToRemove {
			if !isValidNetwork(network) {
				msg := fmt.Sprintf("Network [%v] for account [%v] is "+
					"not a valid configuration value. It must be a "+
					"UUID, CIDR or the string 'public'", network, account)
				log.Fatal(msg)
			}
		}
	}
}

// isValidNetwork validates that a given "network" is specified as expected.
// The expectation is that a "network" is a UUID, CIDR address or the
// string literal 'public'.
func isValidNetwork(network string) bool {
	_, uuidErr := uuid.Parse(network)
	if uuidErr == nil {
		return true
	}

	_, _, ipErr := net.ParseCIDR(network)
	if ipErr == nil {
		return true
	}

	return network == "public"
}
