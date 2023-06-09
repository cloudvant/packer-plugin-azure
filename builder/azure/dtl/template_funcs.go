// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dtl

import (
	"bytes"
	"text/template"

	packertpl "github.com/hashicorp/packer-plugin-sdk/template"
)

func isValidByteValue(b byte) bool {
	if '0' <= b && b <= '9' {
		return true
	}
	if 'a' <= b && b <= 'z' {
		return true
	}
	if 'A' <= b && b <= 'Z' {
		return true
	}
	return b == '.' || b == '_' || b == '-'
}

// Clean up image name by replacing invalid characters with "-"
// Names are not allowed to end in '.', '-', or  '_' and are trimmed.
func templateCleanImageName(s string) string {
	if ok, _ := assertManagedImageName(s, ""); ok {
		return s
	}
	b := []byte(s)
	newb := make([]byte, len(b))
	for i := range newb {
		if isValidByteValue(b[i]) {
			newb[i] = b[i]
		} else {
			newb[i] = '-'
		}
	}

	newb = bytes.TrimRight(newb, "-_.")
	return string(newb)
}

var TemplateFuncs = template.FuncMap{
	"clean_resource_name": templateCleanImageName,
	"clean_image_name":    packertpl.DeprecatedTemplateFunc("clean_image_name", "clean_resource_name", templateCleanImageName),
}
