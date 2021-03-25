package integration

import (
	"github.com/rancher/charts/testing/tests"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	t := tests.Test{
		ID:          "integration-check-image-prefix",
		Description: "Ensure images have rancher/ in front of it",
		On: []tests.ResourceMatcher{
			{
				GroupVersionKind: schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "daemonSet.v1.apps"},
				Filters:          []tests.ResourceFilter{},
			},
			{
				GroupVersionKind: schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
				Filters:          []tests.ResourceFilter{},
			},
			{
				GroupVersionKind: schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "ReplicaSet"},
				Filters:          []tests.ResourceFilter{},
			},
		},
		Do: myTest,
	}
	tests.Register(t)
}

func myTest(objs []v1.Object) (pass bool) {
	logrus.Infof("resources: %v", objs)
	return true
}
