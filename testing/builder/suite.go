package builder

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
)

type testSuite struct {
	// prefix will be added onto the name of all tests that are run
	prefix  string
	decoder runtime.Decoder

	templateObjs map[string][]lintcontext.Object
	tests        []*testBuilder
}

func NewTestSuite(prefix string) *testSuite {
	return &testSuite{
		prefix:       prefix,
		decoder:      serializer.NewCodecFactory(scheme.Scheme).UniversalDeserializer(),
		templateObjs: map[string][]lintcontext.Object{},
		tests:        []*testBuilder{},
	}
}

type ParseOptions struct {
	lintcontext.Options

	// Strict fails to parse a template if a single object cannot be decoded
	// Otherwise, the default behavior is just to print a warning log
	Strict bool
}

func (s *testSuite) ParseTemplate(template, glob string) error {
	return s.ParseTemplateWithOptions(template, glob, ParseOptions{})
}

func (s *testSuite) ParseTemplateWithOptions(template, glob string, options ParseOptions) error {
	if _, exists := s.templateObjs[template]; exists {
		return fmt.Errorf("Cannot parse template %s twice", glob)
	}
	if options.CustomDecoder == nil {
		options.CustomDecoder = s.decoder
	}
	objs, err := parseTemplate(glob, options)
	if err != nil {
		return err
	}
	s.templateObjs[template] = objs
	return nil
}

func (s *testSuite) Test() *testBuilder {
	templates := make([]string, len(s.templateObjs))
	i := 0
	for template := range s.templateObjs {
		templates[i] = template
		i++
	}
	return s.TestTemplates(templates)
}

func (s *testSuite) TestTemplates(templates []string) (b *testBuilder) {
	testTemplateObjs := make(map[string][]lintcontext.Object, len(templates))
	for _, template := range templates {
		objs, exists := s.templateObjs[template]
		if !exists {
			panic(fmt.Sprintf("Template %s has not been added to this suite", template))
		}
		testTemplateObjs[template] = objs
	}
	test := &testBuilder{
		templateObjs: testTemplateObjs,
		prefix:       s.prefix,
	}
	s.tests = append(s.tests, test)
	return test
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
	s.decoder = serializer.NewCodecFactory(customScheme).UniversalDeserializer()
}

func (s *testSuite) SetCustomDecoder(decoder runtime.Decoder) {
	s.decoder = decoder
}
