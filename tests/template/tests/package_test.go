package tests

import "k8s.io/apimachinery/pkg/runtime"

// NOTE: This file is called by main_test.go to initialize package-specific variables.
//       Once generated, this file will never be automatically updated.

const (
	packageName = "rancher-pushprox"
)

func addPackageSchemes(s *runtime.Scheme) {
	// If your package has CRDs, call AddToScheme here to load the CRDs into the deserializer
	// e.g. monitoringv1.AddToScheme(s)
}
