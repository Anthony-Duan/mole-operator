mod-fresh: mod-on mod-tidy mod-vendor mod-off

mod-off:
	export GO111MODULE=off

mod-on:
	export GO111MODULE=on

mod-tidy:
	go mod tidy

mod-vendor:
	go mod vendor

sdk-gen-k8s:
	operator-sdk generate k8s

sdk-gen-crds:
	operator-sdk generate crds

build-cmd:
	go build cmd/manager/main.go

operator-regist:
	kubectl create -f deploy/crds/operator.dtstack.com_moles_crd.yaml