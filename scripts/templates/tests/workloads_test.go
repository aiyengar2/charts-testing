package tests

import (
	"testing"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
)

func TestWorkloads(t *testing.T) {
	test := suite.Test().
		Name("check-image-prefix-0").
		Description("Ensure images have rancher/ in front of it").
		On("apps/v1", "DaemonSet").
		On("apps/v1", "Deployment").
		On("apps/v1", "ReplicaSet").
		Do(GetResources)
	if err := test.Run(); err != nil {
		t.Fatal(err)
	}
}

type Workloads struct {
	Deployments []appsv1.Deployment
	DaemonSets  []appsv1.DaemonSet
	ReplicaSets []appsv1.ReplicaSet
}

func GetResources(w Workloads) (pass bool) {
	logrus.Infof("resources: %v", w)
	return true
}
