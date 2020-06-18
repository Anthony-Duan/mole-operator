package model

const (
	MoleServiceAccountName     = "dtstack"
	MoleServiceName            = "mole-service"
	MoleConfigName             = "mole-config"
	MoleIngressName            = "mole-ingress"
	MoleDeploymentName         = "mole-deployment"
	MolePodName                = "mole-pod"
	MoleHealthEndpoint         = "/api/health"
	MoleHttpPort           int = 3000
	MoleHttpPortName           = "mole-port"
	MoleConfigMountPath        = "/config"
	MoleConfigVolumeName       = "mole-mount"
)
