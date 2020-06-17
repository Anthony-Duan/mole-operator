package model

import (
	molev1 "dtstack.com/dtstack/mole-operator/pkg/apis/mole/v1"
	"github.com/integr8ly/grafana-operator/v3/pkg/apis/integreatly/v1alpha1"
	"k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getIngressTLS(cr *molev1.Mole, name string) []v1beta1.IngressTLS {
	if cr.Spec.Product.Service[name].Instance.Ingress == nil {
		return nil
	}

	if cr.Spec.Product.Service[name].Instance.Ingress.TLSEnabled {
		return []v1beta1.IngressTLS{
			{
				Hosts:      []string{cr.Spec.Product.Service[name].Instance.Ingress.Hostname},
				SecretName: cr.Spec.Product.Service[name].Instance.Ingress.TLSSecretName,
			},
		}
	}
	return nil
}

func getIngressSpec(cr *molev1.Mole, name string) v1beta1.IngressSpec {
	return v1beta1.IngressSpec{
		TLS: getIngressTLS(cr, name),
		Rules: []v1beta1.IngressRule{
			{
				Host: GetHost(cr, name),
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							{
								Path: GetPath(cr, name),
								Backend: v1beta1.IngressBackend{
									ServiceName: BuildResourceName(MoleIngressName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
									ServicePort: GetIngressTargetPort(cr, name),
								},
							},
						},
					},
				},
			},
		},
	}
}

func MoleIngress(cr *molev1.Mole, name string) *v1beta1.Ingress {
	return &v1beta1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name:        BuildResourceName(MoleIngressName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			Namespace:   cr.Namespace,
			Labels:      GetIngressLabels(cr, name),
			Annotations: GetIngressAnnotations(cr, nil, name),
		},
		Spec: getIngressSpec(cr, name),
	}
}

func MoleIngressSelector(cr *molev1.Mole, name string) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      BuildResourceLabel(cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
	}
}

func GetIngressLabels(cr *molev1.Mole, name string) map[string]string {
	if cr.Spec.Product.Service[name].Instance.Ingress == nil {
		return nil
	}
	return cr.Spec.Product.Service[name].Instance.Ingress.Labels
}

func GetIngressAnnotations(cr *molev1.Mole, existing map[string]string, name string) map[string]string {
	if cr.Spec.Product.Service[name].Instance.Ingress == nil {
		return existing
	}
	return MergeAnnotations(cr.Spec.Product.Service[name].Instance.Ingress.Annotations, existing)
}

func GetHost(cr *molev1.Mole, name string) string {
	if cr.Spec.Product.Service[name].Instance.Ingress == nil {
		return ""
	}
	return cr.Spec.Product.Service[name].Instance.Ingress.Hostname
}

func GetPath(cr *molev1.Mole, name string) string {
	if cr.Spec.Product.Service[name].Instance.Ingress == nil {
		return "/"
	}
	return cr.Spec.Product.Service[name].Instance.Ingress.Path
}

func GetIngressTargetPort(cr *molev1.Mole, name string) intstr.IntOrString {
	defaultPort := intstr.FromInt(GetMolePort(cr, name))

	if cr.Spec.Product.Service[name].Instance.Ingress == nil {
		return defaultPort
	}

	if cr.Spec.Product.Service[name].Instance.Ingress.TargetPort == "" {
		return defaultPort
	}

	return intstr.FromString(cr.Spec.Product.Service[name].Instance.Ingress.TargetPort)
}