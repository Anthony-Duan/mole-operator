package apis

import (
	v1 "gitlab.prod.dtstack.cn/dt-insight-ops/mole-operator/pkg/apis/mole/v1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1.SchemeBuilder.AddToScheme)
}
