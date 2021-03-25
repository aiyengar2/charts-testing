package integration

import (
	"github.com/rancher/charts/testing/tests"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	t := tests.Test{
		ID:          "integration-check-image-prefix",
		Description: "Ensure images have rancher/ in front of it",
		On: []tests.ResourceMatcher{
			{
				GroupVersionKind: schema.FromAPIVersionAndKind("apps/v1", "DaemonSet"),
				// Filters:          []tests.ResourceFilter{},
			},
			{
				GroupVersionKind: schema.FromAPIVersionAndKind("apps/v1", "Deployment"),
				// Filters:          []tests.ResourceFilter{},
			},
			{
				GroupVersionKind: schema.FromAPIVersionAndKind("apps/v1", "ReplicaSet"),
				// Filters:          []tests.ResourceFilter{},
			},
		},
		Do: tests.WrapFunc(myTest),
	}
	tests.Register(t)
}

type Workloads struct {
	Deployments []appsv1.Deployment
	DaemonSets  []appsv1.DaemonSet
	ReplicaSets []appsv1.ReplicaSet
}

func myTest(w Workloads) (pass bool) {
	logrus.Infof("resources: %v", w)
	return true
}
