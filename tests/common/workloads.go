package common

import (
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type Workloads struct {
	Deployments []appsv1.Deployment
	DaemonSets  []appsv1.DaemonSet
	ReplicaSets []appsv1.ReplicaSet
	Pods        []corev1.Pod
}

func CheckRancherImagePrefix(w Workloads) (pass bool) {
	result := true

	for _, deployment := range w.Deployments {
		for _, container := range deployment.Spec.Template.Spec.Containers {
			if !ValidateRancherPrefixed(container.Image) {
				logrus.Errorf("Deployment with image missing rancher prefix: %s", deployment.Name)
				result = false
			}
		}
	}
	for _, daemonset := range w.DaemonSets {
		for _, container := range daemonset.Spec.Template.Spec.Containers {
			if !ValidateRancherPrefixed(container.Image) {
				logrus.Errorf("DaemonSet with image missing rancher prefix: %s", daemonset.Name)
				result = false
			}
		}
	}
	for _, replicaset := range w.ReplicaSets {
		for _, container := range replicaset.Spec.Template.Spec.Containers {
			if !ValidateRancherPrefixed(container.Image) {
				logrus.Errorf("ReplicaSet with image missing rancher prefix: %s", replicaset.Name)
				result = false
			}
		}
	}
	for _, pod := range w.Pods {
		for _, container := range pod.Spec.Containers {
			if !ValidateRancherPrefixed(container.Image) {
				logrus.Errorf("Pod with image missing rancher prefix: %s", pod.Name)
				result = false
			}
		}
	}
	return result
}

func CheckImageExists(w Workloads) (pass bool) {
	result := true
	for _, deployment := range w.Deployments {
		for _, container := range deployment.Spec.Template.Spec.Containers {
			if !ValidateImageExists(container.Image) {
				logrus.Errorf("Image does not exist in DockerHub: %s", container.Image)
				result = false
			}
		}
	}
	for _, daemonset := range w.DaemonSets {
		for _, container := range daemonset.Spec.Template.Spec.Containers {
			if !ValidateImageExists(container.Image) {
				logrus.Errorf("Image does not exist in DockerHub: %s", container.Image)
				result = false
			}
		}
	}
	for _, replicaset := range w.DaemonSets {
		for _, container := range replicaset.Spec.Template.Spec.Containers {
			if !ValidateImageExists(container.Image) {
				logrus.Errorf("Image does not exist in DockerHub: %s", container.Image)
				result = false
			}
		}
	}
	for _, pod := range w.Pods {
		for _, container := range pod.Spec.Containers {
			if !ValidateImageExists(container.Image) {
				logrus.Errorf("Image does not exist in DockerHub: %s", container.Image)
				result = false
			}
		}
	}
	return result
}

func CheckSystemDefaultRegistry(w Workloads) (pass bool) {
	registry := "myfakereg.com"
	result := true
	for _, deployment := range w.Deployments {
		for _, container := range deployment.Spec.Template.Spec.Containers {
			if !ValidateSystemDefaultRegistryAndRancher(container.Image, registry) {
				logrus.Errorf("Deployment with image missing rancher prefix: %s", deployment.Name)
				result = false
			}
		}
	}
	for _, daemonset := range w.DaemonSets {
		for _, container := range daemonset.Spec.Template.Spec.Containers {
			if !ValidateSystemDefaultRegistryAndRancher(container.Image, registry) {
				logrus.Errorf("DaemonSet with image missing rancher prefix: %s", daemonset.Name)
				result = false
			}
		}
	}
	for _, replicaset := range w.DaemonSets {
		for _, container := range replicaset.Spec.Template.Spec.Containers {
			if !ValidateSystemDefaultRegistryAndRancher(container.Image, registry) {
				logrus.Errorf("ReplicaSet with image missing rancher prefix: %s", replicaset.Name)
				result = false
			}
		}
	}
	for _, pod := range w.Pods {
		for _, container := range pod.Spec.Containers {
			if !ValidateSystemDefaultRegistryAndRancher(container.Image, registry) {
				logrus.Errorf("Pod with image missing rancher prefix: %s", pod.Name)
				result = false
			}
		}
	}
	return result
}
