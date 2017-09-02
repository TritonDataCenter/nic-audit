/*
 * Copyright (c) 2017, Joyent, Inc. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
package main

import (
	"log"
)

import (
	"github.com/pborman/getopt"
)

// main is the entry point to the application.
func main() {
	configFile := parseCLIFlagsForConfigFilePath()

	log.Println("NIC Compliance Auditing Tool")
	log.Println("https://github.com/joyent/nic-audit\n")
	log.Printf("Reading configuration from: %v\n", configFile)

	config, configErr := readConfigFromFile(configFile)
	validateConfiguration(config)

	if configErr != nil {
		log.Fatalf("Error reading configuration. Details: %v\n", configErr)
	}

	for i := 0; i < len(config.Accounts); i++ {
		account := config.Accounts[i]
		auditErr := auditAccount(account, config.NicGroups, config)

		if auditErr != nil {
			log.Printf("ERROR: %v", auditErr)
		}
	}
}

// parseCLIFlagsForConfigFilePath parses the command line options
// to determine the path to the required configuration file.
func parseCLIFlagsForConfigFilePath() string {
	configPart := getopt.StringLong("config", 'c',
		"/etc/nic-audit.json5",
		"Path to JSON5 format configuration file")

	getopt.Parse()

	if len(*configPart) < 1 {
		log.Fatal("Configuration file must be specified")
	}

	return *configPart
}
