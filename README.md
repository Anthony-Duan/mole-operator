
[operator-sdk](https://github.com/operator-framework/operator-sdk/blob/master/README.md)
# What
A mole is a small animal with black fur that lives underground.
mole-operator is an operator which is responsible for transforming em2.0 schema to kubernetes workload.
mole-operator like a mole under kubernetes, which can be any service operator depended on schema given.


# Setup && Build
## MacOS
brew install operator-sdk
make mod-fresh [可选，默认带vendor]
make build-cmd

# Operator Sdk操作
## Add api
operator-sdk add api --api-version=operator.dtstack.com/v1 --kind=Mole

Create a new API, under an existing project. This command must be run from the project root directory.
Go Example:
$ operator-sdk add api --api-version=operator.dtstack.com/v1 --kind=Mole
This will scaffold the Mole resource API under pkg/apis/operator/v1/...

## Add Controller
operator-sdk add controller --api-version=operator.dtstack.com/v1 --kind=Mole

This will scaffold a new Controller implementation under pkg/controller/mole/....