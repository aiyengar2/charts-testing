package builder

import (
	"context"
	"fmt"

	"github.com/rancher/charts/testing/tests"
	"github.com/sirupsen/logrus"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type testBuilder struct {
	suite *testSuite
	test  tests.Test

	prefix        string
	templateObjs  map[string][]lintcontext.Object
	ignoreFilters bool
}

func (b *testBuilder) All() *testBuilder {
	b.templateObjs = make(map[string][]lintcontext.Object, len(b.suite.templateObjs))
	for k, v := range b.suite.templateObjs {
		b.templateObjs[k] = v
	}
	return b
}

func (b *testBuilder) Include(template string) *testBuilder {
	if _, exists := b.templateObjs[template]; !exists {
		objs, exists := b.suite.templateObjs[template]
		if !exists {
			panic(fmt.Sprintf("Template %s has not been added to this suite", template))
		}
		b.templateObjs[template] = objs
	}
	return b
}

func (b *testBuilder) Exclude(template string) *testBuilder {
	if _, exists := b.templateObjs[template]; exists {
		delete(b.templateObjs, template)
	}
	return b
}

func (b *testBuilder) Name(name string) *testBuilder {
	b.test.ID = fmt.Sprintf("%s-%s", b.prefix, name)
	return b
}

func (b *testBuilder) Description(description string) *testBuilder {
	b.test.Description = description
	return b
}

func (b *testBuilder) Do(testFunc interface{}) *testBuilder {
	b.test.Do = tests.WrapFunc(testFunc)
	return b
}

func (b *testBuilder) DoRaw(doFunc tests.DoFunc) *testBuilder {
	b.test.Do = doFunc
	return b
}

func (b *testBuilder) On(apiVersion string, kind string) *testBuilder {
	return b.OnFilter(apiVersion, kind, []tests.ResourceFilter{})
}

func (b *testBuilder) OnFilter(apiVersion string, kind string, filters []tests.ResourceFilter) *testBuilder {
	b.test.On = append(b.test.On, tests.ResourceMatcher{
		GroupVersionKind: schema.FromAPIVersionAndKind(apiVersion, kind),
		Filters:          filters,
	})
	return b
}

func (b *testBuilder) OnMatcher(matcher tests.ResourceMatcher) *testBuilder {
	b.test.On = append(b.test.On, matcher)
	return b
}

func (b *testBuilder) OnAll() *testBuilder {
	b.ignoreFilters = true
	return b
}

func (b *testBuilder) Run() error {
	return b.RunWithContext(context.TODO())
}

func (b *testBuilder) RunWithContext(ctx context.Context) error {
	if err := b.validate(); err != nil {
		return err
	}
	failedTemplates := []string{}
	for template, lintObjs := range b.templateObjs {
		logrus.Infof("Running test %s on template %s", b.test.ID, template)
		objs := lintObjs
		if !b.ignoreFilters {
			var err error
			objs, err = b.test.FilterObjects(lintObjs)
			if err != nil {
				return fmt.Errorf("Unable to get resources for test %s: %s", b.test.ID, err)
			}
		}
		if !b.test.Run(ctx, objs) {
			failedTemplates = append(failedTemplates, template)
		}
	}
	if len(failedTemplates) > 0 {
		return fmt.Errorf("test %s failed on the following templates: %v", b.test.ID, failedTemplates)
	}
	return nil
}

func (b *testBuilder) validate() error {
	if len(b.test.ID) == 0 {
		return fmt.Errorf("Test must have a name")
	}
	if len(b.test.On) == 0 && !b.ignoreFilters {
		return fmt.Errorf("No resources to run tests on")
	}
	if b.test.Do == nil {
		return fmt.Errorf("No function provided to execute test")
	}
	if len(b.templateObjs) == 0 {
		return fmt.Errorf("Could not find any templates to run %s on", b.test.ID)
	}
	return nil
}
