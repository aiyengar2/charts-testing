package builder

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/rancher/charts/testing/kubelinter/lintcontext"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

type testSuite struct {
	// prefix will be added onto the name of all tests that are run
	prefix string
	scheme *runtime.Scheme

	templateObjs map[string][]lintcontext.Object
	tests        []*testBuilder
}

func NewTestSuite(prefix string) *testSuite {
	return &testSuite{
		prefix:       prefix,
		scheme:       scheme.Scheme,
		templateObjs: map[string][]lintcontext.Object{},
		tests:        []*testBuilder{},
	}
}

func (s *testSuite) ParseTemplate(glob string) error {
	if _, exists := s.templateObjs[glob]; exists {
		return fmt.Errorf("Cannot parse template %s twice", glob)
	}
	objs, err := parseTemplate(glob, s.scheme, false)
	if err != nil {
		return err
	}
	s.templateObjs[glob] = objs
	return nil
}

func (s *testSuite) ParseTemplateStrict(glob string) error {
	if _, exists := s.templateObjs[glob]; exists {
		return fmt.Errorf("Cannot parse template %s twice", glob)
	}
	objs, err := parseTemplate(glob, s.scheme, true)
	if err != nil {
		return err
	}
	s.templateObjs[glob] = objs
	return nil
}

func (s *testSuite) Test() *testBuilder {
	templates := make([]string, len(s.templateObjs))
	i := 0
	for template := range s.templateObjs {
		templates[i] = template
		i++
	}
	return s.test(templates)
}

func (s *testSuite) TestTemplate(template string) (b *testBuilder) {
	if _, exists := s.templateObjs[template]; !exists {
		panic(fmt.Sprintf("Template %s has not been added to this suite", template))
	}
	return s.test([]string{template})
}

func (s *testSuite) RunTests() error {
	var testErrs *multierror.Error
	for _, test := range s.tests {
		if err := test.Run(); err != nil {
			testErrs = multierror.Append(testErrs, err)
		}
	}
	return testErrs
}

func (s *testSuite) SetCustomScheme(customScheme *runtime.Scheme) {
	s.scheme = customScheme
}

func (s *testSuite) test(templates []string) (b *testBuilder) {
	testTemplateObjs := make(map[string][]lintcontext.Object, len(templates))
	for _, template := range templates {
		testTemplateObjs[template] = s.templateObjs[template]
	}
	test := &testBuilder{
		templateObjs: testTemplateObjs,
		prefix:       s.prefix,
	}
	s.tests = append(s.tests, test)
	return test
}