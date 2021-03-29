package tests

import (
	"testing"

	"github.com/rancher/charts/tests/common"
)

// NOTE: This file is provided as an example of how to write a test using the test suite

func TestSingleVersionPerGroupKind(t *testing.T) {
	test := suite.Test().All().
		Name("check-single-version-per-group-kind").
		Description("Ensure that resources of a particular group-kind do not have multiple versions").
		OnAll().
		Do(common.CheckSingleVersionPerGroupKind)
	if err := test.Run(); err != nil {
		t.Fatal(err)
	}
}
