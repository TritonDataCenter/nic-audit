/*
 * Copyright (c) 2017, Joyent, Inc. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
package main

// deleteByValue deletes a map element by the specified value.
func deleteByValue(m map[string]string, value interface{}) {
	for k, v := range m {
		if value == v {
			delete(m, k)
		}
	}
}
