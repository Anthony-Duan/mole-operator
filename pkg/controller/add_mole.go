package controller

import (
	"gitlab.prod.dtstack.cn/dt-insight-ops/mole-operator/pkg/controller/mole"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, mole.Add)
}
