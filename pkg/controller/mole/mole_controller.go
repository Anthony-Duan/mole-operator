package mole

import (
	"container/list"
	"context"
	molev1 "dtstack.com/dtstack/mole-operator/pkg/apis/mole/v1"
	"dtstack.com/dtstack/mole-operator/pkg/controller/common"
	"fmt"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v1beta12 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
)

var log = logf.Log.WithName("controller_mole")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Mole Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return &ReconcileMole{
		client:  mgr.GetClient(),
		scheme:  mgr.GetScheme(),
		context: ctx,
		cancel:  cancel,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mole-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Mole
	err = c.Watch(&source.Kind{Type: &molev1.Mole{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	//Watch for changes to secondary resource and requeue the owner Mole
	if err = watchSecondaryResource(c, &v12.Deployment{}); err != nil {
		return err
	}

	if err = watchSecondaryResource(c, &v1beta12.Ingress{}); err != nil {
		return err
	}

	if err = watchSecondaryResource(c, &v1.ConfigMap{}); err != nil {
		return err
	}

	if err = watchSecondaryResource(c, &v1.Service{}); err != nil {
		return err
	}

	if err = watchSecondaryResource(c, &v1.ServiceAccount{}); err != nil {
		return err
	}

	return nil
}

func watchSecondaryResource(c controller.Controller, resource runtime.Object) error {
	return c.Watch(&source.Kind{Type: resource}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &molev1.Mole{},
	})
}

// blank assignment to verify that ReconcileMole implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMole{}

// ReconcileMole reconciles a Mole object
type ReconcileMole struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client  client.Client
	scheme  *runtime.Scheme
	context context.Context
	cancel  context.CancelFunc
	//recorder record.EventRecorder
}

// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMole) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Mole")

	// Fetch the Mole instance
	instance := &molev1.Mole{}
	err := r.client.Get(r.context, request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	cr := instance.DeepCopy()

	depends := make(map[string][]string)
	for serviceName := range cr.Spec.Product.Service {
		depends[serviceName] = cr.Spec.Product.Service[serviceName].DependsOn
	}

	deploySeq, err := TopologySort(depends) // deploy sequence by service depends on
	if err != nil {
		return reconcile.Result{}, err
	}

	for _, serviceName := range deploySeq {
		//read current state
		currentState := common.NewServiceState(serviceName)
		err = currentState.Read(r.context, cr, r.client)
		if err != nil {
			log.Error(err, "error reading state")
			return r.manageError(cr, err)
		}

		//get desired status
		reconciler := NewMoleReconciler(serviceName)
		desiredState := reconciler.Reconcile(currentState, cr)

		//run action to achieve desired status
		actionRunner := common.NewServiceActionRunner(r.context, r.client, r.scheme, cr)
		err = actionRunner.RunAll(desiredState)
		if err != nil {
			return r.manageError(cr, err)
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileMole) manageError(cr *molev1.Mole, issue error) (reconcile.Result, error) {
	//r.recorder.Event(cr, "Warning", "ProcessingError", issue.Error())
	cr.Status.Phase = molev1.PhaseFailing
	cr.Status.Message = issue.Error()

	err := r.client.Status().Update(r.context, cr)
	if err != nil {
		// Ignore conflicts, resource might just be outdated.
		if errors.IsConflict(err) {
			err = nil
		}
		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: time.Second * 10}, nil
}

func TopologySort(depends map[string][]string) ([]string, error) {
	queue := list.New()
	result := make([]string, 0)
	names := make([]string, 0)
	count := make(map[string]int)

	for name, dependList := range depends { // init topo
		count[name] = len(dependList)
		names = append(names, name)
	}

	for _, name := range names { // add no depends service in queue
		if len(depends[name]) == 0 {
			queue.PushBack(name)
		}
	}
	for queue.Len() > 0 {
		top := queue.Front()
		queue.Remove(top)
		serviceName := top.Value.(string)
		result = append(result, serviceName)
		for _, depend := range depends[serviceName] {
			count[depend]--
			if count[depend] == 0 { // add no depends service in queue
				queue.PushBack(depend)
			}
		}
	}
	if len(result) < len(names) {
		return nil, fmt.Errorf("can't deploy product on this depends")
	}
	return result, nil
}
