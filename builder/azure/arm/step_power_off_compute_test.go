// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arm

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/packer-plugin-azure/builder/azure/common/constants"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

func TestStepPowerOffComputeShouldFailIfPowerOffFails(t *testing.T) {
	var testSubject = &StepPowerOffCompute{
		powerOff: func(context.Context, string, string) error { return fmt.Errorf("!! Unit Test FAIL !!") },
		say:      func(message string) {},
		error:    func(e error) {},
	}

	stateBag := createTestStateBagStepPowerOffCompute()

	var result = testSubject.Run(context.Background(), stateBag)
	if result != multistep.ActionHalt {
		t.Fatalf("Expected the step to return 'ActionHalt', but got '%d'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == false {
		t.Fatalf("Expected the step to set stateBag['%s'], but it was not.", constants.Error)
	}
}

func TestStepPowerOffComputeShouldPassIfPowerOffPasses(t *testing.T) {
	var testSubject = &StepPowerOffCompute{
		powerOff: func(context.Context, string, string) error { return nil },
		say:      func(message string) {},
		error:    func(e error) {},
	}

	stateBag := createTestStateBagStepPowerOffCompute()

	var result = testSubject.Run(context.Background(), stateBag)
	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%d'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == true {
		t.Fatalf("Expected the step to not set stateBag['%s'], but it was.", constants.Error)
	}
}

func TestStepPowerOffComputeShouldTakeStepArgumentsFromStateBag(t *testing.T) {
	var actualResourceGroupName string
	var actualComputeName string

	var testSubject = &StepPowerOffCompute{
		powerOff: func(ctx context.Context, resourceGroupName string, computeName string) error {
			actualResourceGroupName = resourceGroupName
			actualComputeName = computeName

			return nil
		},
		say:   func(message string) {},
		error: func(e error) {},
	}

	stateBag := createTestStateBagStepPowerOffCompute()
	var result = testSubject.Run(context.Background(), stateBag)

	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%d'.", result)
	}

	var expectedComputeName = stateBag.Get(constants.ArmComputeName).(string)
	var expectedResourceGroupName = stateBag.Get(constants.ArmResourceGroupName).(string)

	if actualComputeName != expectedComputeName {
		t.Fatal("Expected the step to source 'constants.ArmResourceGroupName' from the state bag, but it did not.")
	}

	if actualResourceGroupName != expectedResourceGroupName {
		t.Fatal("Expected the step to source 'constants.ArmResourceGroupName' from the state bag, but it did not.")
	}
}

func createTestStateBagStepPowerOffCompute() multistep.StateBag {
	stateBag := new(multistep.BasicStateBag)

	stateBag.Put(constants.ArmComputeName, "Unit Test: ComputeName")
	stateBag.Put(constants.ArmResourceGroupName, "Unit Test: ResourceGroupName")

	return stateBag
}
