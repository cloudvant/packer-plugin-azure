// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/packer-plugin-azure/builder/azure/common/constants"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

func TestProcessStepResultShouldContinueForNonErrors(t *testing.T) {
	stateBag := new(multistep.BasicStateBag)

	code := processStepResult(nil, func(error) { t.Fatal("Should not be called!") }, stateBag)
	if _, ok := stateBag.GetOk(constants.Error); ok {
		t.Errorf("Error was nil, but was still in the state bag.")
	}

	if code != multistep.ActionContinue {
		t.Errorf("Expected ActionContinue(%d), but got=%d", multistep.ActionContinue, code)
	}
}

func TestProcessStepResultShouldHaltOnError(t *testing.T) {
	stateBag := new(multistep.BasicStateBag)
	isSaidError := false

	code := processStepResult(fmt.Errorf("boom"), func(error) { isSaidError = true }, stateBag)
	if _, ok := stateBag.GetOk(constants.Error); !ok {
		t.Errorf("Error was non nil, but was not in the state bag.")
	}

	if !isSaidError {
		t.Errorf("Expected error to be said, but it was not.")
	}

	if code != multistep.ActionHalt {
		t.Errorf("Expected ActionHalt(%d), but got=%d", multistep.ActionHalt, code)
	}
}
