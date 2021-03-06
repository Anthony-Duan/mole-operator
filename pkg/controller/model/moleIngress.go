package model

import (
	molev1 "gitlab.prod.dtstack.cn/dt-insight-ops/mole-operator/pkg/apis/mole/v1"
	"fmt"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
)

func getIngressTLS(cr *molev1.Mole, name string) []v1beta1.IngressTLS {
	if cr.Spec.Product.Service[name].Instance.Ingress == nil {
		return nil
	}

	if cr.Spec.Product.Service[name].Instance.Ingress.TLSEnabled {
		return []v1beta1.IngressTLS{
			{
				Hosts:      []string{GetHost(cr, name)},
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
						Paths: getIngressRule(cr, name),
					},
				},
			}, {
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: getIngressRulePaths(cr, name),
					},
				},
			},
		},
	}
}

func getIngressRule(cr *molev1.Mole, name string) []v1beta1.HTTPIngressPath {
	paths := make([]v1beta1.HTTPIngressPath, 0)
	for _, port := range cr.Spec.Product.Service[name].Instance.Deployment.Ports {
		paths = append(paths, v1beta1.HTTPIngressPath{
			//Path: fmt.Sprintf("/%v/%v/", cr.Namespace, strings.ToLower(name)),
			Backend: v1beta1.IngressBackend{
				ServiceName: BuildResourceName(MoleServiceName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
				ServicePort: intstr.FromInt(port),
			},
		})
	}
	return paths
}

func getIngressRulePaths(cr *molev1.Mole, name string) []v1beta1.HTTPIngressPath {
	paths := make([]v1beta1.HTTPIngressPath, 0)
	for _, port := range cr.Spec.Product.Service[name].Instance.Deployment.Ports {
		paths = append(paths, v1beta1.HTTPIngressPath{
			Path: fmt.Sprintf("/%v/%v/", cr.Namespace, strings.ToLower(name)),
			Backend: v1beta1.IngressBackend{
				ServiceName: BuildResourceName(MoleServiceName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
				ServicePort: intstr.FromInt(port),
			},
		})
	}
	return paths
}

func MoleIngress(cr *molev1.Mole, name string) *v1beta1.Ingress {
	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        BuildResourceName(MoleIngressName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
			Namespace:   cr.Namespace,
			Labels:      GetIngressLabels(cr, name),
			Annotations: GetIngressAnnotations(cr, nil, name),
		},
		Spec: getIngressSpec(cr, name),
	}
}

func MoleIngressReconciled(cr *molev1.Mole, currentState *v1beta1.Ingress, name string) *v1beta1.Ingress {
	reconciled := currentState.DeepCopy()
	reconciled.Labels = GetIngressLabels(cr, name)
	reconciled.Annotations = GetIngressAnnotations(cr, currentState.Annotations, name)
	reconciled.Spec = getIngressSpec(cr, name)
	return reconciled
}

func MoleIngressSelector(cr *molev1.Mole, name string) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      BuildResourceName(MoleIngressName, cr.Spec.Product.ParentProductName, cr.Spec.Product.ProductName, name),
	}
}

func GetIngressLabels(cr *molev1.Mole, name string) map[string]string {
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

func GetIngressAnnotations(cr *molev1.Mole, existing map[string]string, name string) map[string]string {
	return map[string]string{
		"kubernetes.io/ingress.class":                "nginx",
		"nginx.ingress.kubernetes.io/ssl-redirect":   "false",
		"nginx.ingress.kubernetes.io/rewrite-target": "/",
	}
	//if cr.Spec.Product.Service[name].Instance.Ingress == nil {
	//	return existing
	//}
	//return MergeAnnotations(cr.Spec.Product.Service[name].Instance.Ingress.Annotations, existing)
}

func GetHost(cr *molev1.Mole, name string) string {
	if cr.Spec.Product.Service[name].Instance.Ingress == nil {
		fmt.Sprintf("%v.%v.dtstack.com")
		return fmt.Sprintf("%v.%v.dtstack.com", strings.ToLower(cr.Spec.Product.ProductName), cr.Namespace)
	}
	return cr.Spec.Product.Service[name].Instance.Ingress.Hostname
}

func GetPath(cr *molev1.Mole, name string) string {
	if cr.Spec.Product.Service[name].Instance.Ingress == nil {
		return "/"
	}
	return cr.Spec.Product.Service[name].Instance.Ingress.Path
}

//func GetIngressTargetPort(cr *molev1.Mole, name string) intstr.IntOrString {
//	defaultPort := intstr.FromInt(GetMolePort(cr, name))
//
//	if cr.Spec.Product.Service[name].Instance.Ingress == nil {
//		return defaultPort
//	}
//
//	if cr.Spec.Product.Service[name].Instance.Ingress.TargetPort == "" {
//		return defaultPort
//	}
//
//	return intstr.FromString(cr.Spec.Product.Service[name].Instance.Ingress.TargetPort)
//}
