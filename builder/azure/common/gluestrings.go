// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

// removes overlap between the end of a and the start of b and
// glues them together
func GlueStrings(a, b string) string {
	shift := 0
	for shift < len(a) {
		i := 0
		for (i+shift < len(a)) && (i < len(b)) && (a[i+shift] == b[i]) {
			i++
		}
		if i+shift == len(a) {
			break
		}
		shift++
	}

	return a[:shift] + b
}
