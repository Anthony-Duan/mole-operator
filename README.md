
[operator-sdk](https://github.com/operator-framework/operator-sdk/blob/master/README.md)


# API操作
## Add api
operator-sdk add api --api-version=operator.dtstack.com/v1 --kind=Mole

Create a new API, under an existing project. This command must be run from the project root directory.
Go Example:
$ operator-sdk add api --api-version=operator.dtstack.com/v1 --kind=Mole
This will scaffold the Mole resource API under pkg/apis/operator/v1/...

## Add Controller
operator-sdk add controller --api-version=operator.dtstack.com/v1 --kind=Mole

This will scaffold a new Controller implementation under pkg/controller/mole/....