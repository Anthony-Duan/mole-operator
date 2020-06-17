package model

const (
	MoleServiceAccountName           = "dtstack"
	MoleServiceName                  = "mole-service"
	MoleDataStorageName              = "mole-pvc"
	MoleConfigName                   = "mole-config"
	MoleConfigFileName               = "mole.ini"
	MoleIngressName                  = "mole-ingress"
	MoleRouteName                    = "mole-route"
	MoleDeploymentName               = "mole-deployment"
	MolePodName                      = "mole-pod"
	MolePluginsVolumeName            = "mole-plugins"
	MoleInitContainerName            = "mole-plugins-init"
	MoleLogsVolumeName               = "mole-logs"
	MoleDataVolumeName               = "mole-data"
	MoleDatasourcesConfigMapName     = "mole-datasources"
	MoleHealthEndpoint               = "/api/health"
	MolePodLabel                     = "mole"
	LastConfigAnnotation             = "last-config"
	LastConfigEnvVar                 = "LAST_CONFIG"
	LastDatasourcesConfigEnvVar      = "LAST_DATASOURCES"
	MoleAdminSecretName              = "mole-admin-credentials"
	MoleHttpPort                 int = 3000
	MoleHttpPortName                 = "mole-port"
)
