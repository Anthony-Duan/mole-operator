package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type ConfigMap map[string]string

type ServiceConfig struct {
	ServiceDisplay string `json:"service_display,omitempty"`
	Version        string `json:"version,omitempty"`
	//Instance       struct {
	//	ConfigPaths    []string           `json:"config_paths"`
	//	Logs           []string           `json:"logs"`
	//	DataDir        []string           `json:"data_dir"`
	//	Environment    map[string]string  `json:"environment"`
	//	Cmd            string             `json:"cmd,omitempty"`
	//	PrometheusPort string             `json:"prometheus_port,omitempty"`
	//	Replica        string             `json:"replica,omitempty"`
	//} `json:"instance,omitempty"`
	//Group     string    `json:"group,omitempty"`
	//DependsOn []string  `json:"depends_on"`
	//Config    ConfigMap `json:"config,omitempty"`

	BaseProduct   string `json:"base_product,omitempty"`
	BaseService   string `json:"base_service,omitempty"`
	BaseParsed    bool   `json:"base_parsed,omitempty"`
	BaseAtrribute string `json:"base_atrribute,omitempty"`
}

type SchemaConfig struct {
	ParentProductName  string                   `json:"parent_product_name,omitempty"`
	ProductName        string                   `json:"product_name,omitempty"`
	ProductNameDisplay string                   `json:"product_name_display,omitempty"`
	ProductVersion     string                   `json:"product_version,omitempty"`
	Service            map[string]ServiceConfig `json:"service"`
}

// MoleSpec defines the desired state of Mole
type MoleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Product SchemaConfig `json:"product,omitempty"`
}

// MoleStatus defines the observed state of Mole
type MoleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Mole is the Schema for the moles API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=moles,scope=Namespaced
type Mole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MoleSpec   `json:"spec,omitempty"`
	Status MoleStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MoleList contains a list of Mole
type MoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Mole{}, &MoleList{})
}
