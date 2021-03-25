package tests

import (
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type ResourceFilter struct {
	// Resourceis the metadata.name of the object
	Name string `json:"name" yaml:"name"`
	// Namespace is the metadata.namespace of the object
	Namespace string `json:"namespace" yaml:"namespace"`
	// LabelSelector is the k8s label selector used to match against this object
	LabelSelector v1.LabelSelector `json:"labelSelector" yaml:"labelSelector"`
	// Annotations are metadata.annotations keys that must exist on the object
	Annotations []string `json:"annotations" yaml:"annotations"`
}

func (f ResourceFilter) String() string {
	return fmt.Sprintf("{name: %s, namespace: %v, labelSelector: %s, annotation: %s}", f.Name, f.Namespace, &f.LabelSelector, f.Annotations)
}

func (f *ResourceFilter) Match(lintObj lintcontext.Object) (bool, error) {
	obj := lintObj.K8sObject
	if len(f.Name) > 0 && obj.GetName() != f.Name {
		return false, nil
	}
	if len(f.Namespace) > 0 && obj.GetNamespace() != f.Namespace {
		return false, nil
	}
	// Check labels
	selector, err := v1.LabelSelectorAsSelector(&f.LabelSelector)
	if err != nil {
		return false, fmt.Errorf("Failed to parse labelSelector for %v: %s", f, err)
	}
	if !selector.Matches(labels.Set(obj.GetLabels())) {
		return false, nil
	}
	// Check annotations
	for _, annotation := range f.Annotations {
		hasAnnotation := false
		for key := range obj.GetAnnotations() {
			if key == annotation {
				hasAnnotation = true
				break
			}
		}
		if !hasAnnotation {
			return false, nil
		}
	}
	return true, nil
}
