package tests

import (
	"fmt"

	"github.com/rancher/charts/testing/kubelinter/lintcontext"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// func Match(resourceType )

// ResourceMatcher determines whether a resource should be matched
type ResourceMatcher struct {
	schema.GroupVersionKind `json:"gvk" yaml:"gvk"`

	Filters []ResourceFilter `json:"filters" yaml:"filters"`
}

func (m ResourceMatcher) String() string {
	return fmt.Sprintf("{gvk: %s, filters: %v}", m.GroupVersionKind, m.Filters)
}

func (m *ResourceMatcher) Match(lintObj lintcontext.Object) (bool, error) {
	obj := lintObj.K8sObject
	if obj.GetObjectKind().GroupVersionKind() != m.GroupVersionKind {
		return false, nil
	}
	if len(m.Filters) == 0 {
		return true, nil
	}
	for _, f := range m.Filters {
		match, err := f.Match(lintObj)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}
