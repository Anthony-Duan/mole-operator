package model

import (
	molev1 "dtstack.com/dtstack/mole-operator/pkg/apis/mole/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getServiceLabels(cr *molev1.Mole, name string) map[string]string {
	if cr.Spec.Product.Service[name].Instance.Service == nil {
		return nil
	}
	return cr.Spec.Product.Service[name].Instance.Service.Labels
}

func getServiceAnnotations(cr *molev1.Mole, existing map[string]string, name string) map[string]string {
	if cr.Spec.Product.Service[name].Instance.Service.Annotations == nil {
		return existing
	}

	return MergeAnnotations(cr.Spec.Product.Service[name].Instance.Service.Annotations, existing)
}

func getServiceType(cr *molev1.Mole, name string) v1.ServiceType {
	if cr.Spec.Product.Service[name].Instance.Service == nil {
		return v1.ServiceTypeClusterIP
	}
	if cr.Spec.Product.Service[name].Instance.Service.Type == "" {
		return v1.ServiceTypeClusterIP
	}
	return cr.Spec.Product.Service[name].Instance.Service.Type
}

func GetMolePort(cr *molev1.Mole, name string) int {
	return cr.Spec.Product.Service[name].Instance.ContainerPort
}

func getServicePorts(cr *molev1.Mole, currentState *v1.Service, name string) []v1.ServicePort {
	ContainerPort := int32(GetMolePort(cr, name))
	portName := BuildPortName(cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, cr.Spec.Product.ProductVersion, name, MoleHttpPortName)
	defaultPorts := []v1.ServicePort{
		{
			Name:       portName,
			Protocol:   "TCP",
			Port:       ContainerPort,
			TargetPort: intstr.FromString(portName),
		},
	}

	if cr.Spec.Product.Service[name].Instance.Service == nil {
		return defaultPorts
	}

	return defaultPorts
}

func MoleService(cr *molev1.Mole, name string) *v1.Service {
	return &v1.Service{
		ObjectMeta: v12.ObjectMeta{
			Name:        BuildResourceName(MoleServiceName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			Namespace:   cr.Namespace,
			Labels:      getServiceLabels(cr, name),
			Annotations: getServiceAnnotations(cr, nil, name),
		},
		Spec: v1.ServiceSpec{
			Ports: getServicePorts(cr, nil, name),
			Selector: map[string]string{
				"app": BuildResourceLabel(cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			},
			ClusterIP: "",
			Type:      getServiceType(cr, name),
		},
	}
}

func MoleServiceReconciled(cr *molev1.Mole, currentState *v1.Service, name string) *v1.Service {
	reconciled := currentState.DeepCopy()
	reconciled.Labels = getServiceLabels(cr, name)
	reconciled.Annotations = getServiceAnnotations(cr, currentState.Annotations, name)
	reconciled.Spec.Ports = getServicePorts(cr, currentState, name)
	reconciled.Spec.Type = getServiceType(cr, name)
	return reconciled
}

func MoleServiceSelector(cr *molev1.Mole, name string) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      BuildResourceLabel(cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
	}
}
