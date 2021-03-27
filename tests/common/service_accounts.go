package common

import (
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
)

type ServiceAccountUsages struct {
	ServiceAccounts          []corev1.ServiceAccount
	Deployments              []appsv1.Deployment
	Jobs                     []batchv1.Job
	ClusterRoleBindingsBeta1 []rbacv1beta1.ClusterRoleBinding
	ClusterRoleBindingsV1    []rbacv1.ClusterRoleBinding
	RoleBindings             []rbacv1.RoleBinding
}

func CheckServiceAccountUsage(s ServiceAccountUsages) (pass bool) {
	unusedServiceAccounts := make(map[string]bool)
	for _, sa := range s.ServiceAccounts {
		unusedServiceAccounts[sa.Name] = true
	}

	for _, deployment := range s.Deployments {
		if _, ok := unusedServiceAccounts[deployment.Spec.Template.Spec.ServiceAccountName]; ok {
			delete(unusedServiceAccounts, deployment.Spec.Template.Spec.ServiceAccountName)
		}
	}
	for _, job := range s.Jobs {
		if _, ok := unusedServiceAccounts[job.Spec.Template.Spec.ServiceAccountName]; ok {
			delete(unusedServiceAccounts, job.Spec.Template.Spec.ServiceAccountName)
		}
	}
	for _, crbs := range s.ClusterRoleBindingsBeta1 {
		for _, subject := range crbs.Subjects {
			if subject.Kind == "ServiceAccount" {
				if _, ok := unusedServiceAccounts[subject.Name]; ok {
					delete(unusedServiceAccounts, subject.Name)
				}
			}
		}
	}
	for _, crbs := range s.ClusterRoleBindingsV1 {
		for _, subject := range crbs.Subjects {
			if subject.Kind == "ServiceAccount" {
				if _, ok := unusedServiceAccounts[subject.Name]; ok {
					delete(unusedServiceAccounts, subject.Name)
				}
			}
		}
	}
	for _, rbs := range s.RoleBindings {
		for _, subject := range rbs.Subjects {
			if subject.Kind == "ServiceAccount" {
				if _, ok := unusedServiceAccounts[subject.Name]; ok {
					delete(unusedServiceAccounts, subject.Name)
				}
			}
		}
	}
	if len(unusedServiceAccounts) == 0 {
		return true
	} else {
		logrus.Errorf("Unused Service Accounts: %v\n", unusedServiceAccounts)
		return false
	}
}
