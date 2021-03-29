package tests

import (
	"context"
	"fmt"

	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Want to initialize a way to get a specific set of resources
// e.g. Get me all Deployments + a specific ServiceAccount + a specific PodSecurityPolicy
// Once I ensure that those resources I expect exist, I want to be able to look through them via a type I unmarshall them into

type Test struct {
	// ID is a unique identifier for this test
	ID string `json:"id" yaml:"id"`
	// Description describes what is being tested
	Description string `json:"description" yaml:"description"`
	// On contains a list of ResourceMatchers that match against resources that meet the criteria
	On []ResourceMatcher `json:"on" yaml:"on"`
	// Do is the test that will be called on all resources that were returned by ResourceGetters
	Do DoFunc `json:"do" yaml:"do"`
}

func (t Test) String() string {
	return fmt.Sprintf("{test: %s, on: %v, description: %s}", t.ID, t.On, t.Description)
}

func (t *Test) Run(ctx context.Context, lintObjs []lintcontext.Object) (pass bool) {
	objs := make([]v1.Object, len(lintObjs))
	for i, lintObj := range lintObjs {
		objs[i] = lintObj.K8sObject
	}
	return t.Do(ctx, objs)
}

func (t *Test) FilterObjects(lintObjs []lintcontext.Object) (objs []lintcontext.Object, err error) {
	for _, lintObj := range lintObjs {
		for _, resourceMatcher := range t.On {
			match, err := resourceMatcher.Match(lintObj)
			if err != nil {
				return nil, err
			}
			if match {
				objs = append(objs, lintObj)
			}
		}
	}
	return objs, nil
}
