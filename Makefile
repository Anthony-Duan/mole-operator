mod: mod-on mod-tidy mod-vendor mod-off

mod-off:
	export GO111MODULE=off

mod-on:
	export GO111MODULE=on

mod-tidy:
	go mod tidy

mod-vendor:
	go mod vendor



crds:
    # operator-sdk generate crds --crd-version v1beta1
	operator-sdk generate crds

deepcopy:
	operator-sdk generate k8s

image:
	operator-sdk build registry.cn-hangzhou.aliyuncs.com/dtstack/mole:v1.0.17
    docker tag registry.cn-hangzhou.aliyuncs.com/dtstack/mole:v1.0.17 172.16.8.120:5443/dtstack-dev/mole:v1.0.17

push:
	docker push registry.cn-hangzhou.aliyuncs.com/dtstack/mole:v1.0.17
	docker push 172.16.8.120:5443/dtstack-dev/mole:v1.0.17

gobuild:
	go build cmd/manager/main.go

run:
	operator-sdk run local --watch-namespace="dtstack-system"

registcrd:
	kubectl create -f deploy/crds/operator.dtstack.com_moles_crd.yaml
