package lintcontext

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

func UseCustomScheme(s *runtime.Scheme) {
	clientSchema = s
	decoder = serializer.NewCodecFactory(s).UniversalDeserializer()
}
