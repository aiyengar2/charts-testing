package common

import (
	"context"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/runtime"
)

type GenericObject interface {
	v1.Object
	runtime.Object
}

type GenericObjects struct {
	Objects []GenericObject
}

func CheckSingleVersionPerGroupKind(ctx context.Context, objs GenericObjects) (pass bool) {
	groupKindMap := map[schema.GroupKind]string{}
	for _, obj := range objs.Objects {
		gvk := obj.GetObjectKind().GroupVersionKind()
		gk := gvk.GroupKind()
		if version, exists := groupKindMap[gk]; exists && version != gvk.Version {
			logrus.Errorf("Found conflicting versions for resources of GroupKind %s: %s", gk, []string{gvk.Version, version})
			return false
		}
	}
	return true
}
