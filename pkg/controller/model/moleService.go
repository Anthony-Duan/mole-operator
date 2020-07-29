package model

import (
	molev1 "dtstack.com/dtstack/mole-operator/pkg/apis/mole/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

func getServiceLabels(cr *molev1.Mole, name string) map[string]string {
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

func getServiceAnnotations(cr *molev1.Mole, existing map[string]string, name string) map[string]string {
	if cr.Spec.Product.Service[name].Instance.Service == nil || cr.Spec.Product.Service[name].Instance.Service.Annotations == nil {
		return existing
	}
	return MergeAnnotations(cr.Spec.Product.Service[name].Instance.Service.Annotations, existing)
}

func getServiceType(cr *molev1.Mole, name string) corev1.ServiceType {
	if cr.Spec.Product.Service[name].Instance.Service == nil {
		return corev1.ServiceTypeClusterIP
	}
	if cr.Spec.Product.Service[name].Instance.Service.Type == "" {
		return corev1.ServiceTypeClusterIP
	}
	return cr.Spec.Product.Service[name].Instance.Service.Type
}

//func GetMolePort(cr *molev1.Mole, name string, index int) int {
//	return cr.Spec.Product.Service[name].Instance.Deployment.Ports[index]
//}

func getServicePorts(cr *molev1.Mole, name string) []corev1.ServicePort {
	//portName := BuildPortName(name, MoleHttpPortName)
	defaultPorts := make([]corev1.ServicePort, 0)
	for index, port := range cr.Spec.Product.Service[name].Instance.Deployment.Ports {
		defaultPorts = append(defaultPorts, corev1.ServicePort{
			Name:       BuildPortName(name, index),
			Protocol:   "TCP",
			Port:       int32(port),
			TargetPort: intstr.FromString(BuildPortName(name, index)),
		})
	}
	if cr.Spec.Product.Service[name].Instance.Service == nil {
		return defaultPorts
	}
	return defaultPorts
}

func MoleService(cr *molev1.Mole, name string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        BuildResourceName(MoleServiceName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			Namespace:   cr.Namespace,
			Labels:      getServiceLabels(cr, name),
			Annotations: getServiceAnnotations(cr, nil, name),
		},
		Spec: corev1.ServiceSpec{
			Ports: getServicePorts(cr, name),
			Selector: map[string]string{
				"app": BuildResourceLabel(cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			},
			ClusterIP: "",
			Type:      getServiceType(cr, name),
		},
	}
}

func MoleServiceReconciled(cr *molev1.Mole, currentState *corev1.Service, name string) *corev1.Service {
	reconciled := currentState.DeepCopy()
	reconciled.Labels = getServiceLabels(cr, name)
	reconciled.Annotations = getServiceAnnotations(cr, currentState.Annotations, name)
	reconciled.Spec.Ports = getServicePorts(cr, name)
	reconciled.Spec.Type = getServiceType(cr, name)
	return reconciled
}

func MoleServiceSelector(cr *molev1.Mole, name string) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      BuildResourceName(MoleServiceName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
	}
}
