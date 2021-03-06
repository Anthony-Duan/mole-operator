package model

import (
	molev1 "gitlab.prod.dtstack.cn/dt-insight-ops/mole-operator/pkg/apis/mole/v1"
	"strconv"
)

func GetMoleLabels(cr *molev1.Mole) map[string]string {
	var labels = map[string]string{}

	labels["pid"] = strconv.Itoa(cr.Spec.Product.Pid)
	labels["deploy_uuid"] = cr.Spec.Product.DeployUUid
	labels["cluster_id"] = strconv.Itoa(cr.Spec.Product.ClusterId)
	labels["product_name"] = cr.Spec.Product.ProductName
	labels["product_version"] = cr.Spec.Product.ProductVersion
	labels["parent_product_name"] = cr.Spec.Product.ParentProductName
	labels["com"] = MoleCom

	return labels
}
