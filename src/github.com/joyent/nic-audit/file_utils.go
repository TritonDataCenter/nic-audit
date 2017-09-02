/*
 * Copyright (c) 2017, Joyent, Inc. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
package main

import (
	"golang.org/x/sys/unix"
	"os"
)

// exists function that determines if a given path exists.
func exists(filePath string) (exists bool) {
	exists = true

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		exists = false
	}

	return exists
}

// isReadable determines if a given directory or file can be read from.
func isReadable(path string) (readable bool) {
	return unix.Access(path, unix.R_OK) == nil
}
