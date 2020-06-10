package controller

import (
	"dtstack.com/dtstack/mole-operator/pkg/controller/mole"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, mole.Add)
}
