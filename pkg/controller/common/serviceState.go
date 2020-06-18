package common

import (
	"context"
	molev1 "dtstack.com/dtstack/mole-operator/pkg/apis/mole/v1"
	"dtstack.com/dtstack/mole-operator/pkg/controller/model"
	v13 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceState struct {
	Name           string
	MoleConfig     *v1.ConfigMap
	MoleIngress    *v1beta1.Ingress
	MoleService    *v1.Service
	MoleDeployment *v13.Deployment
}

func NewServiceState(name string) *ServiceState {
	return &ServiceState{
		Name: name,
	}
}

func (i *ServiceState) Read(ctx context.Context, cr *molev1.Mole, client client.Client) error {
	err := i.readMoleDeployment(ctx, cr, client, i.Name)
	if err != nil {
		return err
	}

	err = i.readMoleService(ctx, cr, client, i.Name)
	if err != nil {
		return err
	}

	err = i.readMoleIngress(ctx, cr, client)
	return err
}

func (i *ServiceState) readMoleService(ctx context.Context, cr *molev1.Mole, client client.Client, name string) error {
	currentState := model.MoleService(cr, name)
	selector := model.MoleServiceSelector(cr, name)
	err := client.Get(ctx, selector, currentState)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	i.MoleService = currentState.DeepCopy()
	return nil
}

func (i *ServiceState) readMoleIngress(ctx context.Context, cr *molev1.Mole, client client.Client) error {
	currentState := model.MoleIngress(cr, i.Name)
	selector := model.MoleIngressSelector(cr, i.Name)
	err := client.Get(ctx, selector, currentState)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	i.MoleIngress = currentState.DeepCopy()
	return nil
}

func (i *ServiceState) readMoleDeployment(ctx context.Context, cr *molev1.Mole, client client.Client, name string) error {
	currentState := model.MoleDeployment(cr, name)
	selector := model.MoleDeploymentSelector(cr)
	err := client.Get(ctx, selector, currentState)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	i.MoleDeployment = currentState.DeepCopy()
	return nil
}
