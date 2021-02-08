/*
Copyright 2021 The OpDev Team.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	corev1 "k8s.io/api/core/v1"

	"github.com/opdev/aws-account-binding-operator/controllers/states"
	. "github.com/opdev/aws-account-binding-operator/helpers/reconcileresults"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	awsint "github.com/opdev/aws-account-binding-operator/api/v1alpha1"
)

// AWSAccountBindingReconciler reconciles a AWSAccountBinding object
type AWSAccountBindingReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=integrations.opdev.io,resources=awsaccountbindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=integrations.opdev.io,resources=awsaccountbindings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=integrations.opdev.io,resources=awsaccountbindings/finalizers,verbs=update

// Reconcile will attempt to make the cluster state of a given
// AWSAccountBinding match the desired state.

// Reconcile handles events indicating that the desired state of AWSAccountBinding resources
// may have changed, and does what's necessary to make the existing state match the desired state.
func (r *AWSAccountBindingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = context.WithValue(ctx, instKeyContextKey, req.NamespacedName)
	_ = r.Log.WithValues("AWSAccountBinding", req.NamespacedName)
	r.Log.Info(fmt.Sprintf("starting reconciliation for %s", req.NamespacedName))
	defer r.Log.Info(fmt.Sprintf("ending reconciliation for %s", req.NamespacedName))

	state, res, err := r.DetermineState(ctx)
	if ShouldHaltOrRequeue(res, err) {
		r.Log.Info("Reconcile() halting while calling DetermineState")
		return Evaluate(res, err)
	}

	// handle deletion
	deletionSubReconcilers := []subreconcilerFuncs{
		r.removeNamespaceAnnotation,
		r.removeFinalizer,
	}

	if state.IsBeingDeleted() {
		r.Log.Info("resource is being deleted, running deletion reconciliation flows")
		for _, f := range deletionSubReconcilers {
			if r, err := f(ctx); ShouldHaltOrRequeue(r, err) {
				return Evaluate(r, err)
			}
		}

		return Evaluate(DoNotRequeue())
	}

	r.Log.Info("running reconciliation flows")
	subreconcilers := []subreconcilerFuncs{
		r.addFinalizer,
		r.updateStatus,
		r.addNamespaceAnnotation,
	}

	for _, f := range subreconcilers {
		// call the reconciler with the state
		if r, err := f(ctx); ShouldHaltOrRequeue(r, err) {
			return Evaluate(r, err)
		}
	}

	// successfully reconciled
	return Evaluate(DoNotRequeue())
}

// GetInstance queries the API for the instance of the custom resource.
func (r *AWSAccountBindingReconciler) GetInstance(ctx context.Context) (awsint.AWSAccountBinding, *ctrl.Result, error) {
	instanceKey := ctx.Value(instKeyContextKey).(types.NamespacedName)
	var instance awsint.AWSAccountBinding
	if err := r.Get(ctx, instanceKey, &instance); err != nil {
		if apierrors.IsNotFound(err) {
			// it was deleted before reconcile completed
			r.Log.Info("GetInstance() resource not found, it was likely deleted.")
			cres, e := DoNotRequeue()
			return awsint.AWSAccountBinding{}, cres, e
		}

		cres, e := RequeueWithError(err)
		return awsint.AWSAccountBinding{}, cres, e
	}

	cres, e := ContinueReconciling()
	return instance, cres, e
}

// GetNamespace queries the API for the namespace associated with the custom resource request.
func (r *AWSAccountBindingReconciler) GetNamespace(ctx context.Context) (corev1.Namespace, *ctrl.Result, error) {
	instanceKey := ctx.Value(instKeyContextKey).(types.NamespacedName)

	ns := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: instanceKey.Name},
	}

	// the namespace should exist, but if it doesn't we cannot continue
	if err := r.Get(ctx, client.ObjectKeyFromObject(&ns), &ns); err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Error(err, "unable to continue with reconciliation if associated namespace does not exist")
			// do not requeue because we don't want to cause a loop
			cres, e := DoNotRequeue()
			return corev1.Namespace{}, cres, e
		}
		cres, e := RequeueWithError(err)
		return corev1.Namespace{}, cres, e
	}

	cres, e := ContinueReconciling()
	return ns, cres, e
}

// GetResources queries the API for resources necessary to determine the state
// of the existing AWSAccountBinding
func (r *AWSAccountBindingReconciler) GetResources(ctx context.Context) (states.AccountBindingResources, *ctrl.Result, error) {
	instance, res, err := r.GetInstance(ctx)
	if ShouldHaltOrRequeue(res, err) {
		r.Log.Info("GetResources() halting while calling GetInstance")
		return states.AccountBindingResources{}, res, err
	}

	ns, res, err := r.GetNamespace(ctx)
	if ShouldHaltOrRequeue(res, err) {
		r.Log.Info("GetResources() halting while calling GetNamespace")
		return states.AccountBindingResources{}, res, err
	}

	r.Log.Info("GetResources() completed successfully")
	cres, e := ContinueReconciling()
	return states.NewAccountBindingResources(instance, ns), cres, e
}

// DetermineState queries the API for resources necessary to determine the state
// of existing resources, and then returns the state.
func (r *AWSAccountBindingReconciler) DetermineState(ctx context.Context) (states.AccountBindingState, *ctrl.Result, error) {
	resource, res, err := r.GetResources(ctx)
	if ShouldHaltOrRequeue(res, err) {
		r.Log.Info("DetermineState() halting while calling GetResources")
		return states.AccountBindingState{}, res, err
	}

	r.Log.Info("DetermineState() completed successfully")
	cres, e := ContinueReconciling()
	return resource.ParseState(), cres, e
}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSAccountBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&awsint.AWSAccountBinding{}).
		Owns(&corev1.Namespace{}).
		// WithEventFilter(predicate.Funcs{
		// 	UpdateFunc:  func(ue event.UpdateEvent) bool { return false },
		// 	GenericFunc: func(ge event.GenericEvent) bool { return false },
		// 	CreateFunc:  func(ce event.CreateEvent) bool { return true },
		// 	DeleteFunc:  func(de event.DeleteEvent) bool { return true },
		// }).
		// TODO build a predicate or use a predicate that prevents updates
		// or controller restarts from re-reconciling an account binding.
		Complete(r)
}
