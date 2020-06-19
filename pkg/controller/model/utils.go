package model

import (
	"fmt"
)

func MergeAnnotations(requested map[string]string, existing map[string]string) map[string]string {
	//if existing == nil {
	//	return requested
	//}
	//
	//for k, v := range requested {
	//	existing[k] = v
	//}
	return existing
}

func BuildResourceName(resourceType, parentProductName, productName, serviceName string) string {
	return fmt.Sprintf("%v-%v-%v-%v", resourceType, parentProductName, productName, serviceName)
}

func BuildResourceLabel(parentProductName, productName, serviceName string) string {
	return fmt.Sprintf("%v-%v-%v", parentProductName, productName, serviceName)
}

func BuildConfigMapName(parentProductName, productName, productVersion, serviceName, configName string) string {
	return fmt.Sprintf("%v-%v-%v-%v-%v", parentProductName, productName, productVersion, serviceName, configName)
}
func BuildPortName(serviceName, portName string) string {
	return fmt.Sprintf("%v-%v", serviceName, portName)
}
