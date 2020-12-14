package controller

import (
	"github.com/Anthony-Duan/mole-operator/pkg/controller/mole"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, mole.Add)
}
