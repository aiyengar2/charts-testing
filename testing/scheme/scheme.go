package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubectl/pkg/scheme"
)

func Scheme() *runtime.Scheme {
	baseScheme := scheme.Scheme
	return AddCRDsToScheme(baseScheme)
}

func AddCRDsToScheme(s *runtime.Scheme) *runtime.Scheme {
	return s
}
