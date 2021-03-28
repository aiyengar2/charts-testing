package tests

import (
	"testing"

	"github.com/rancher/charts/tests/common"
)

// NOTE: This file is provided as an example of how to write a test using the test suite

func TestAllServiceAccountsAreUsed(t *testing.T) {
	test := suite.Test().All().
		Name("check-service-account-usage").
		Description("Ensure each service account is used at least once").
		On("apps/v1", "Deployment").
		On("batch/v1", "Job").
		On("v1", "ServiceAccount").
		On("rbac.authorization.k8s.io/v1beta1", "ClusterRoleBinding").
		On("rbac.authorization.k8s.io/v1", "ClusterRoleBinding").
		Do(common.CheckServiceAccountUsage)
	if err := test.Run(); err != nil {
		t.Fatal(err)
	}
}
