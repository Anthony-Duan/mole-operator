package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MolePhase string

var (
	MOLEF_RUNNING MolePhase = "Running"
	MOLE_PENDING  MolePhase = "Pending"
	MOLE_FAILED   MolePhase = "Failed"
)

type ConfigMap map[string]string

type Instance struct {
	ConfigPaths    []string                 `json:"config_path,omitempty"`
	Logs           []string                 `json:"logs,omitempty"`
	DataDir        []string                 `json:"data_dir,omitempty"`
	Environment    map[string]string        `json:"environment,omitempty"`
	Cmd            string                   `json:"cmd,omitempty,omitempty"`
	PrometheusPort string                   `json:"prometheus_port,omitempty"`
	Ingress        *MoleIngress             `json:"ingress,omitempty"`
	Service        *MoleService             `json:"service,omitempty"`
	Deployment     *MoleDeployment          `json:"deployment,omitempty"`
	Resources      *v1.ResourceRequirements `json:"resources,omitempty"`
	PostDeploy     string                   `json:"post_deploy,omitempty"`
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
	IsJob           bool     `json:"is_job,omitempty"`
}

type SchemaConfig struct {
	Pid                int                      `json:"pid,omitempty"`
	ClusterId          int                      `json:"cluster_id,omitempty"`
	DeployUUid         string                   `json:"deploy_uuid,omitempty"`
	ParentProductName  string                   `json:"parent_product_name,omitempty"`
	ProductName        string                   `json:"product_name,omitempty"`
	ProductNameDisplay string                   `json:"product_name_display,omitempty"`
	ProductVersion     string                   `json:"product_version,omitempty"`
	ImagePullSecret    string                   `json:"imagePullSecret,omitempty"`
	Service            map[string]ServiceConfig `json:"service"`
}

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

type MoleService struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Type        v1.ServiceType    `json:"type,omitempty"`
	Ports       []v1.ServicePort  `json:"ports,omitempty"`
}

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

type MoleSpec struct {
	Product SchemaConfig `json:"product,omitempty"`
}

type MoleStatus struct {
	Phase   MolePhase `json:"phase"`
	Message string    `json:"message"`
}

type Mole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MoleSpec   `json:"spec,omitempty"`
	Status MoleStatus `json:"status,omitempty"`
}

type MoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Mole{}, &MoleList{})
}
