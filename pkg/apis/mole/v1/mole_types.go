package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type StatusPhase string

var (
	NoPhase          StatusPhase
	PhaseReconciling StatusPhase = "reconciling"
	PhaseFailing     StatusPhase = "failing"
)

type ConfigMap map[string]string

type Instance struct {
	ConfigPath     string                   `json:"config_path,omitempty"`
	Logs           []string                 `json:"logs,omitempty"`
	DataDir        []string                 `json:"data_dir,omitempty"`
	Environment    map[string]string        `json:"environment,omitempty"`
	Cmd            string                   `json:"cmd,omitempty,omitempty"`
	PrometheusPort string                   `json:"prometheus_port,omitempty"`
	Ingress        *MoleIngress             `json:"ingress,omitempty"`
	Service        *MoleService             `json:"service,omitempty"`
	Deployment     *MoleDeployment          `json:"deployment,omitempty"`
	Resources      *v1.ResourceRequirements `json:"resources,omitempty"`
}

type ServiceConfig struct {
	ServiceDisplay  string   `json:"service_display,omitempty"`
	IsDeployIngress bool     `json:"is_deploy_ingress,omitempty"`
	Version         string   `json:"version,omitempty"`
	Instance        Instance `json:"instance,omitempty"`
	Group           string   `json:"group,omitempty"`
	DependsOn       []string `json:"depends_on,omitempty"`
	BaseProduct     string   `json:"base_product,omitempty"`
	BaseService     string   `json:"base_service,omitempty"`
	BaseParsed      bool     `json:"base_parsed,omitempty"`
	BaseAttribute   string   `json:"base_attribute,omitempty"`
}

type SchemaConfig struct {
	ParentProductName  string                   `json:"parent_product_name,omitempty"`
	ProductName        string                   `json:"product_name,omitempty"`
	ProductNameDisplay string                   `json:"product_name_display,omitempty"`
	ProductVersion     string                   `json:"product_version,omitempty"`
	ProductUUid        string                   `json:"product_uuid,omitempty"`
	Service            map[string]ServiceConfig `json:"service"`
}

// MoleIngress provides a means to configure the ingress created
type MoleIngress struct {
	Annotations   map[string]string `json:"annotations,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	Hostname      string            `json:"hostname,omitempty"`
	Path          string            `json:"path,omitempty"`
	Enabled       bool              `json:"enabled,omitempty"`
	TLSEnabled    bool              `json:"tlsEnabled,omitempty"`
	TLSSecretName string            `json:"tlsSecretName,omitempty"`
	TargetPort    string            `json:"targetPort,omitempty"`
}

// MoleService provides a means to configure the service
type MoleService struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Type        v1.ServiceType    `json:"type,omitempty"`
	Ports       []v1.ServicePort  `json:"ports,omitempty"`
}

// MoleDeployment provides a means to configure the deployment
type MoleDeployment struct {
	Annotations                   map[string]string      `json:"annotations,omitempty"`
	Labels                        map[string]string      `json:"labels,omitempty"`
	Replicas                      int32                  `json:"replicas,omitempty"`
	Image                         string                 `json:"image,omitempty"`
	Ports                         []int                  `json:"ports,omitempty"`
	Containers                    []MoleContainer        `json:"containers,omitempty"`
	NodeSelector                  map[string]string      `json:"nodeSelector,omitempty"`
	Tolerations                   []v1.Toleration        `json:"tolerations,omitempty"`
	Affinity                      *v1.Affinity           `json:"affinity,omitempty"`
	SecurityContext               *v1.PodSecurityContext `json:"securityContext,omitempty"`
	TerminationGracePeriodSeconds int64                  `json:"terminationGracePeriodSeconds,omitempty"`
}

type MoleContainer struct {
	Image string `json:"image,omitempty"`
	Name  string `json:"name,omitempty"`
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
	Phase   StatusPhase `json:"phase"`
	Message string      `json:"message"`
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
