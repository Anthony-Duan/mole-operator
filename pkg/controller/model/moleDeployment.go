package model

import (
	molev1 "dtstack.com/dtstack/mole-operator/pkg/apis/mole/v1"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

func getAffinities(cr *molev1.Mole, name string) *v13.Affinity {
	var affinity = v13.Affinity{}
	if cr.Spec.Product.Service[name].Instance.Deployment != nil && cr.Spec.Product.Service[name].Instance.Deployment.Affinity != nil {
		affinity = *cr.Spec.Product.Service[name].Instance.Deployment.Affinity
	}
	return &affinity
}

func getSecurityContext(cr *molev1.Mole, name string) *v13.PodSecurityContext {
	var securityContext = v13.PodSecurityContext{}
	if cr.Spec.Product.Service[name].Instance.Deployment != nil && cr.Spec.Product.Service[name].Instance.Deployment.SecurityContext != nil {
		securityContext = *cr.Spec.Product.Service[name].Instance.Deployment.SecurityContext
	}
	return &securityContext
}

func getReplicas(cr *molev1.Mole, name string) *int32 {
	var replicas int32 = 1
	if cr.Spec.Product.Service[name].Instance.Deployment == nil {
		return &replicas
	}
	if cr.Spec.Product.Service[name].Instance.Deployment.Replicas <= 0 {
		return &replicas
	} else {
		return &cr.Spec.Product.Service[name].Instance.Deployment.Replicas
	}
}

func getRollingUpdateStrategy() *v1.RollingUpdateDeployment {
	var maxUnaval intstr.IntOrString = intstr.FromInt(25)
	var maxSurge intstr.IntOrString = intstr.FromInt(25)
	return &v1.RollingUpdateDeployment{
		MaxUnavailable: &maxUnaval,
		MaxSurge:       &maxSurge,
	}
}

func getPodAnnotations(cr *molev1.Mole, existing map[string]string, name string) map[string]string {
	var annotations = map[string]string{}
	// Add fixed annotations
	annotations["prometheus.io/scrape"] = "true"
	//annotations["prometheus.io/port"] = fmt.Sprintf("%v", GetMolePort(cr, name))
	annotations = MergeAnnotations(annotations, existing)

	if cr.Spec.Product.Service[name].Instance.Deployment != nil {
		annotations = MergeAnnotations(cr.Spec.Product.Service[name].Instance.Deployment.Annotations, annotations)
	}
	return annotations
}

func getPodLabels(cr *molev1.Mole, name string) map[string]string {
	var labels = map[string]string{}

	labels["app"] = BuildResourceLabel(cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name)
	labels["pid"] = strconv.Itoa(cr.Spec.Product.Pid)
	labels["deploy_uuid"] = cr.Spec.Product.DeployUUid
	labels["cluster_id"] = strconv.Itoa(cr.Spec.Product.ClusterId)
	labels["product_name"] = cr.Spec.Product.ProductName
	labels["product_version"] = cr.Spec.Product.ProductVersion
	labels["parent_product_name"] = cr.Spec.Product.ParentProductName
	labels["service_name"] = name
	labels["service_version"] = cr.Spec.Product.Service[name].Version
	labels["group"] = cr.Spec.Product.Service[name].Group
	labels["com"] = MoleCom

	return labels
}
func getDeploymentLabels(cr *molev1.Mole, name string) map[string]string {
	var labels = map[string]string{}

	labels["pid"] = strconv.Itoa(cr.Spec.Product.Pid)
	labels["deploy_uuid"] = cr.Spec.Product.DeployUUid
	labels["cluster_id"] = strconv.Itoa(cr.Spec.Product.ClusterId)
	labels["product_name"] = cr.Spec.Product.ProductName
	labels["product_version"] = cr.Spec.Product.ProductVersion
	labels["parent_product_name"] = cr.Spec.Product.ParentProductName
	labels["service_name"] = name
	labels["service_version"] = cr.Spec.Product.Service[name].Version
	labels["group"] = cr.Spec.Product.Service[name].Group
	labels["com"] = MoleCom

	return labels
}

func getNodeSelectors(cr *molev1.Mole, name string) map[string]string {
	var nodeSelector = map[string]string{}

	if cr.Spec.Product.Service[name].Instance.Deployment != nil && cr.Spec.Product.Service[name].Instance.Deployment.NodeSelector != nil {
		nodeSelector = cr.Spec.Product.Service[name].Instance.Deployment.NodeSelector
	}
	return nodeSelector

}

func getTerminationGracePeriod(cr *molev1.Mole, name string) *int64 {
	var tcp int64 = 30
	if cr.Spec.Product.Service[name].Instance.Deployment != nil && cr.Spec.Product.Service[name].Instance.Deployment.TerminationGracePeriodSeconds != 0 {
		tcp = cr.Spec.Product.Service[name].Instance.Deployment.TerminationGracePeriodSeconds
	}
	return &tcp

}

func getTolerations(cr *molev1.Mole, name string) []v13.Toleration {
	tolerations := []v13.Toleration{}

	if cr.Spec.Product.Service[name].Instance.Deployment != nil && cr.Spec.Product.Service[name].Instance.Deployment.Tolerations != nil {
		for _, val := range cr.Spec.Product.Service[name].Instance.Deployment.Tolerations {
			tolerations = append(tolerations, val)
		}
	}
	return tolerations
}

func getVolumes(cr *molev1.Mole, name string) []v13.Volume {
	var volumes []v13.Volume
	// Volume to mount the config file from a configMap
	volumes = append(volumes, v13.Volume{
		Name: BuildResourceName(MoleConfigVolumeName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
		VolumeSource: v13.VolumeSource{
			ConfigMap: &v13.ConfigMapVolumeSource{
				LocalObjectReference: v13.LocalObjectReference{
					Name: BuildResourceName(MoleConfigName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
				},
			},
		},
	})

	//Volume to mount emptyDir to share logs
	hostPathType := v13.HostPathDirectoryOrCreate
	volumes = append(volumes, v13.Volume{
		Name: BuildResourceName(MoleLogsVolumeName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
		VolumeSource: v13.VolumeSource{
			HostPath: &v13.HostPathVolumeSource{
				Path: LogHostPath,
				Type: &hostPathType,
			},
		},
	})
	return volumes
}

func getVolumeMounts(cr *molev1.Mole, name string) []v13.VolumeMount {
	var mounts []v13.VolumeMount
	for _, configPath := range cr.Spec.Product.Service[name].Instance.ConfigPaths {
		mounts = append(mounts, v13.VolumeMount{
			Name:      BuildResourceName(MoleConfigVolumeName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			MountPath: configPath,
		})
	}
	mounts = append(mounts, v13.VolumeMount{
		Name:      BuildResourceName(MoleLogsVolumeName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
		MountPath: "/log/data",
	})

	return mounts
}

//func getProbe(cr *molev1.Mole, delay, timeout, failure int32, name string) *v13.Probe {
//	return &v13.Probe{
//		Handler: v13.Handler{
//			HTTPGet: &v13.HTTPGetAction{
//				Path: MoleHealthEndpoint,
//				Port: intstr.FromInt(GetMolePort(cr, name)),
//			},
//		},
//		InitialDelaySeconds: delay,
//		TimeoutSeconds:      timeout,
//		FailureThreshold:    failure,
//	}
//}

func getContainers(cr *molev1.Mole, name string) []v13.Container {
	var containers []v13.Container
	containers = append(containers, v13.Container{
		Name:         name,
		Image:        cr.Spec.Product.Service[name].Instance.Deployment.Image,
		WorkingDir:   "",
		Ports:        getContainerPorts(cr, name),
		VolumeMounts: getVolumeMounts(cr, name),
		//LivenessProbe:            getProbe(cr, 60, 30, 10, name),
		//ReadinessProbe:           getProbe(cr, 5, 3, 1, name),
		//TerminationMessagePath:   "/dev/termination-log",
		//TerminationMessagePolicy: "File",
		ImagePullPolicy: "IfNotPresent",
	})
	for _, container := range cr.Spec.Product.Service[name].Instance.Deployment.Containers {
		containers = append(containers, v13.Container{
			Name:            container.Name,
			Image:           container.Image,
			VolumeMounts:    getVolumeMounts(cr, name),
			ImagePullPolicy: "IfNotPresent",
		})
	}

	return containers
}

func getContainerPorts(cr *molev1.Mole, name string) []v13.ContainerPort {
	//portName := BuildPortName(name, MoleHttpPortName)
	defaultPorts := make([]v13.ContainerPort, 0)
	for index, port := range cr.Spec.Product.Service[name].Instance.Deployment.Ports {
		defaultPorts = append(defaultPorts, v13.ContainerPort{
			Name:          BuildPortName(name, index),
			Protocol:      "TCP",
			ContainerPort: int32(port),
		})
	}
	return defaultPorts
}

func getDeploymentSpec(cr *molev1.Mole, annotations map[string]string, name string) v1.DeploymentSpec {
	return v1.DeploymentSpec{
		Replicas: getReplicas(cr, name),
		Selector: &v12.LabelSelector{
			MatchLabels: map[string]string{
				"app": BuildResourceLabel(cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			},
		},
		Template: v13.PodTemplateSpec{
			ObjectMeta: v12.ObjectMeta{
				Name:        BuildResourceName(MolePodName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
				Labels:      getPodLabels(cr, name),
				Annotations: getPodAnnotations(cr, annotations, name),
			},
			Spec: v13.PodSpec{
				NodeSelector:     getNodeSelectors(cr, name),
				Tolerations:      getTolerations(cr, name),
				Affinity:         getAffinities(cr, name),
				SecurityContext:  getSecurityContext(cr, name),
				Volumes:          getVolumes(cr, name),
				Containers:       getContainers(cr, name),
				ImagePullSecrets: getImagePullSecrets(cr),
				//ServiceAccountName: MoleServiceAccountName,
				//RestartPolicy:   v13.RestartPolicyAlways,
				//TerminationGracePeriodSeconds: getTerminationGracePeriod(cr, name),
			},
		},
		Strategy: v1.DeploymentStrategy{
			Type:          "RollingUpdate",
			RollingUpdate: getRollingUpdateStrategy(),
		},
	}
}

func MoleDeployment(cr *molev1.Mole, name string) *v1.Deployment {
	return &v1.Deployment{
		ObjectMeta: v12.ObjectMeta{
			Name:      BuildResourceName(MoleDeploymentName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			Labels:    getDeploymentLabels(cr, name),
			Namespace: cr.Namespace,
		},
		Spec: getDeploymentSpec(cr, nil, name),
	}
}

func MoleDeploymentReconciled(cr *molev1.Mole, currentState *v1.Deployment, name string) *v1.Deployment {
	reconciled := currentState.DeepCopy()
	reconciled.Labels = getDeploymentLabels(cr, name)
	reconciled.Spec = getDeploymentSpec(cr, currentState.Spec.Template.Annotations, name)
	return reconciled
}

func MoleDeploymentSelector(cr *molev1.Mole, name string) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      BuildResourceName(MoleDeploymentName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
	}
}

func getImagePullSecrets(cr *molev1.Mole) []v13.LocalObjectReference {
	return []v13.LocalObjectReference{
		{Name: cr.Spec.Product.ImagePullSecret},
	}
}
