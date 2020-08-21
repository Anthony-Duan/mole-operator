package model

import (
	molev1 "dtstack.com/dtstack/mole-operator/pkg/apis/mole/v1"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
)

const (
	MemoryRequest = "100Mi"
	MemoryLimit   = "8Gi"
	CpuLimit      = "4000m"
	CpuRequest    = "0"
)

func getAffinities(cr *molev1.Mole, name string) *corev1.Affinity {
	var affinity = corev1.Affinity{}
	if cr.Spec.Product.Service[name].Instance.Deployment != nil && cr.Spec.Product.Service[name].Instance.Deployment.Affinity != nil {
		affinity = *cr.Spec.Product.Service[name].Instance.Deployment.Affinity
	}
	return &affinity
}

func getSecurityContext(cr *molev1.Mole, name string) *corev1.PodSecurityContext {
	var securityContext = corev1.PodSecurityContext{}
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

func getRollingUpdateStrategy() *appsv1.RollingUpdateDeployment {
	var maxUnavailable = intstr.FromInt(0)
	var maxSurge = intstr.FromString("25%")
	return &appsv1.RollingUpdateDeployment{
		MaxUnavailable: &maxUnavailable,
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

func getTolerations(cr *molev1.Mole, name string) []corev1.Toleration {
	tolerations := []corev1.Toleration{}

	if cr.Spec.Product.Service[name].Instance.Deployment != nil && cr.Spec.Product.Service[name].Instance.Deployment.Tolerations != nil {
		for _, val := range cr.Spec.Product.Service[name].Instance.Deployment.Tolerations {
			tolerations = append(tolerations, val)
		}
	}
	return tolerations
}

func getVolumes(cr *molev1.Mole, name string) []corev1.Volume {
	var volumes []corev1.Volume
	// Volume to mount the config file from a configMap
	volumes = append(volumes, corev1.Volume{
		Name: BuildResourceName(MoleConfigVolumeName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: BuildResourceName(MoleConfigName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
				},
			},
		},
	})

	//Volume to mount hostPath to share logs
	//hostPathType := corev1.HostPathDirectoryOrCreate
	//volumes = append(volumes, corev1.Volume{
	//	Name: BuildResourceName(MoleLogsVolumeName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
	//	VolumeSource: corev1.VolumeSource{
	//		HostPath: &corev1.HostPathVolumeSource{
	//			Path: LogPath + "/" + cr.Spec.Product.ProductName + "/" + name,
	//			Type: &hostPathType,
	//		},
	//	},
	//})
	return volumes
}

func getVolumeMounts(cr *molev1.Mole, name string) []corev1.VolumeMount {
	var mounts []corev1.VolumeMount
	for _, configPath := range cr.Spec.Product.Service[name].Instance.ConfigPaths {
		subPath := strings.Replace(configPath, "/", "_", -1)
		mounts = append(mounts, corev1.VolumeMount{
			Name:      BuildResourceName(MoleConfigVolumeName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			SubPath:   subPath,
			MountPath: fmt.Sprintf("opt/dtstack/%v/%v/%v", cr.Spec.Product.ProductName, name, configPath),
		})
	}
	//mounts = append(mounts, corev1.VolumeMount{
	//	Name:      BuildResourceName(MoleLogsVolumeName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
	//	MountPath: MoleMountPath,
	//})

	return mounts
}

//func getProbe(cr *molev1.Mole, delay, timeout, failure int32, name string) *corev1.Probe {
//	return &corev1.Probe{
//		Handler: corev1.Handler{
//			TCPSocket: &corev1.TCPSocketAction{
//				Port: intstr.FromInt(),
//			},
//		},
//		InitialDelaySeconds: delay,
//		TimeoutSeconds:      timeout,
//		FailureThreshold:    failure,
//	}
//}

func getContainers(cr *molev1.Mole, name string) []corev1.Container {
	var containers []corev1.Container
	containers = append(containers, corev1.Container{
		Name:            ConvertDNSRuleName(name),
		Image:           cr.Spec.Product.Service[name].Instance.Deployment.Image,
		WorkingDir:      "",
		Ports:           getContainerPorts(cr, name),
		VolumeMounts:    getVolumeMounts(cr, name),
		Resources:       getResources(cr, name),
		ImagePullPolicy: "Always",
		//Lifecycle:       getPodLifeCycle(),
		//LivenessProbe:  getProbe(cr, 0, 10, 10, name),
		//ReadinessProbe: getProbe(cr, 0, 3, 1, name),
		//TerminationMessagePath:   "/dev/termination-log",
		//TerminationMessagePolicy: "File",
	})
	for _, container := range cr.Spec.Product.Service[name].Instance.Deployment.Containers {
		containers = append(containers, corev1.Container{
			Name:            container.Name,
			Image:           container.Image,
			VolumeMounts:    getVolumeMounts(cr, name),
			ImagePullPolicy: "IfNotPresent",
		})
	}

	return containers
}

func getContainerPorts(cr *molev1.Mole, name string) []corev1.ContainerPort {
	//portName := BuildPortName(name, MoleHttpPortName)
	defaultPorts := make([]corev1.ContainerPort, 0)
	for index, port := range cr.Spec.Product.Service[name].Instance.Deployment.Ports {
		defaultPorts = append(defaultPorts, corev1.ContainerPort{
			Name:          BuildPortName(name, index),
			Protocol:      "TCP",
			ContainerPort: int32(port),
		})
	}
	return defaultPorts
}

func getDeploymentSpec(cr *molev1.Mole, annotations map[string]string, name string) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Replicas:        getReplicas(cr, name),
		MinReadySeconds: 10,

		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": BuildResourceLabel(cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:        BuildResourceName(MolePodName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
				Labels:      getPodLabels(cr, name),
				Annotations: getPodAnnotations(cr, annotations, name),
			},
			Spec: corev1.PodSpec{
				NodeSelector:     getNodeSelectors(cr, name),
				Tolerations:      getTolerations(cr, name),
				Affinity:         getAffinities(cr, name),
				SecurityContext:  getSecurityContext(cr, name),
				Volumes:          getVolumes(cr, name),
				Containers:       getContainers(cr, name),
				ImagePullSecrets: getImagePullSecrets(cr),

				//ServiceAccountName: MoleServiceAccountName,
				//RestartPolicy:   corev1.RestartPolicyAlways,
				//TerminationGracePeriodSeconds: getTerminationGracePeriod(cr, name),
			},
		},
		Strategy: appsv1.DeploymentStrategy{
			Type:          "RollingUpdate",
			RollingUpdate: getRollingUpdateStrategy(),
		},
	}
}

func MoleDeployment(cr *molev1.Mole, name string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BuildResourceName(MoleDeploymentName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			Labels:    getDeploymentLabels(cr, name),
			Namespace: cr.Namespace,
		},
		Spec: getDeploymentSpec(cr, nil, name),
	}
}

func MoleDeploymentReconciled(cr *molev1.Mole, currentState *appsv1.Deployment, name string) *appsv1.Deployment {
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

func getImagePullSecrets(cr *molev1.Mole) []corev1.LocalObjectReference {
	return []corev1.LocalObjectReference{
		{Name: cr.Spec.Product.ImagePullSecret},
	}
}

func getPodLifeCycle() *corev1.Lifecycle {
	return &corev1.Lifecycle{
		PostStart: &corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"/bin/sh",
					"-c",
					"mkdir -p /mount/${HOSTNAME}/logs && ln -s /mount/${HOSTNAME}/logs logs",
				},
			},
		},
	}
}
func getResources(cr *molev1.Mole, name string) corev1.ResourceRequirements {
	if cr.Spec.Product.Service[name].Instance.Deployment.Resources != nil {
		return *cr.Spec.Product.Service[name].Instance.Deployment.Resources
	}
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(MemoryRequest),
			corev1.ResourceCPU:    resource.MustParse(CpuRequest),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(MemoryLimit),
			corev1.ResourceCPU:    resource.MustParse(CpuLimit),
		},
	}
}
