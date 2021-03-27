package tests

import (
	"testing"

	"github.com/rancher/charts/tests/common"
)

func TestRancherImage(t *testing.T) {
	test := suite.Test().
		Name("check-image-prefix").
		Description("Ensure images have rancher/ in front of them").
		On("apps/v1", "DaemonSet").
		On("apps/v1", "Deployment").
		On("apps/v1", "ReplicaSet").
		On("v1", "Pod").
		Do(common.CheckRancherImagePrefix)
	if err := test.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestRancherImageWithPrivateRegistry(t *testing.T) {
	test := suite.Test().
		Name("check-system-default-registry").
		Description("Ensure images use <system-default-registry>/rancher/ if it's provided").
		On("apps/v1", "DaemonSet").
		On("apps/v1", "Deployment").
		On("apps/v1", "ReplicaSet").
		On("v1", "Pod").
		Do(common.CheckSystemDefaultRegistry)
	if err := test.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestImageExists(t *testing.T) {
	test := suite.Test().
		Name("check-image-exists").
		Description("Ensure images are present in DockerHub").
		On("apps/v1", "DaemonSet").
		On("apps/v1", "Deployment").
		On("apps/v1", "ReplicaSet").
		On("v1", "Pod").
		Do(common.CheckImageExists)
	if err := test.Run(); err != nil {
		t.Fatal(err)
	}
}
