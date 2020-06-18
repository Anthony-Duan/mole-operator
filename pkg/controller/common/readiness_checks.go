package common

import (
	"errors"
	v12 "k8s.io/api/apps/v1"
	"k8s.io/api/extensions/v1beta1"
)

const (
	ConditionStatusSuccess = "True"
)

func IsIngressReady(ingress *v1beta1.Ingress) bool {
	if ingress == nil {
		return false
	}

	return len(ingress.Status.LoadBalancer.Ingress) > 0
}

func IsDeploymentReady(deployment *v12.Deployment) (bool, error) {
	if deployment == nil {
		return false, nil
	}
	// A deployment has an array of conditions
	for _, condition := range deployment.Status.Conditions {
		// One failure condition exists, if this exists, return the Reason
		if condition.Type == v12.DeploymentReplicaFailure {
			return false, errors.New(condition.Reason)
			// A successful deployment will have the progressing condition type as true
		} else if condition.Type == v12.DeploymentProgressing && condition.Status != ConditionStatusSuccess {
			return false, nil
		}
	}

	return deployment.Status.ReadyReplicas == deployment.Status.Replicas, nil
}
