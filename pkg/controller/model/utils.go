package model

import (
	"fmt"
	"strings"
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
	serviceName = strings.Replace(serviceName, "_", "", -1)
	productName = strings.Replace(productName, "_", "", -1)
	parentProductName = strings.Replace(parentProductName, "_", "", -1)
	return fmt.Sprintf("%v-%v-%v-%v", resourceType, parentProductName, productName, serviceName)
}

func BuildResourceLabel(parentProductName, productName, serviceName string) string {
	serviceName = strings.Replace(serviceName, "_", "", -1)
	productName = strings.Replace(productName, "_", "", -1)
	parentProductName = strings.Replace(parentProductName, "_", "", -1)
	return fmt.Sprintf("%v-%v-%v", parentProductName, productName, serviceName)
}

func BuildConfigMapName(parentProductName, productName, serviceName, configName string) string {
	serviceName = strings.Replace(serviceName, "_", "", -1)
	productName = strings.Replace(productName, "_", "", -1)
	parentProductName = strings.Replace(parentProductName, "_", "", -1)
	return fmt.Sprintf("%v-%v-%v-%v", parentProductName, productName, serviceName, configName)
}
func BuildPortName(serviceName string, index int) string {
	return fmt.Sprintf("%v-%v", serviceName, index)
}
