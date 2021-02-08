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

	"github.com/go-logr/logr"
	"github.com/opdev/aws-account-binding-operator/controllers/states"
	reconc "github.com/opdev/aws-account-binding-operator/helpers/reconcileresults"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	awsint "github.com/opdev/aws-account-binding-operator/api/v1alpha1"
)

// AWSAccountBindingApprovalReconciler reconciles a AWSAccountBindingApproval object
type AWSAccountBindingApprovalReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=integrations.opdev.io,resources=awsaccountbindingapprovals,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=integrations.opdev.io,resources=awsaccountbindingapprovals/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=integrations.opdev.io,resources=awsaccountbindingapprovals/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *AWSAccountBindingApprovalReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = context.WithValue(ctx, instKeyContextKey, req.NamespacedName)
	_ = r.Log.WithValues("AWSAccountBindingApproval", req.NamespacedName)
	r.Log.Info(fmt.Sprintf("starting reconciliation for %s", req.NamespacedName))
	defer r.Log.Info(fmt.Sprintf("ending reconciliation for %s", req.NamespacedName))

	state, res, err := r.DetermineState(ctx)
	if reconc.ShouldHaltOrRequeue(res, err) {
		r.Log.Info("Reconcile() halting while calling DetermineState")
		return reconc.Evaluate(res, err)
	}

	// handle deletion
	deletionSubReconcilers := []subreconcilerFuncs{
		r.removeAccountBinding,
		r.removeFinalizer,
	}

	if state.IsBeingDeleted() {
		r.Log.Info("resource is being deleted, running deletion reconciliation flows")
		for _, f := range deletionSubReconcilers {
			if r, err := f(ctx); reconc.ShouldHaltOrRequeue(r, err) {
				return reconc.Evaluate(r, err)
			}
		}

		return reconc.Evaluate(reconc.DoNotRequeue())
	}

	r.Log.Info("running reconciliation flows")
	subreconcilers := []subreconcilerFuncs{
		r.updateStatus,
		r.ensureApproval,
		r.addFinalizer,
		r.createAccountBinding,
	}

	for _, f := range subreconcilers {
		// call the reconciler with the state
		if r, err := f(ctx); reconc.ShouldHaltOrRequeue(r, err) {
			return reconc.Evaluate(r, err)
		}
	}

	return reconc.Evaluate(reconc.DoNotRequeue())
}

// GetInstance queries the API for the instance of the custom resource.
func (r *AWSAccountBindingApprovalReconciler) GetInstance(ctx context.Context) (awsint.AWSAccountBindingApproval, *ctrl.Result, error) {
	instanceKey := ctx.Value(instKeyContextKey).(types.NamespacedName)
	var instance awsint.AWSAccountBindingApproval
	if err := r.Get(ctx, instanceKey, &instance); err != nil {
		if apierrors.IsNotFound(err) {
			// it was deleted before reconcile completed
			r.Log.Info("GetInstance() resource not found, it was likely deleted.")
			cres, e := reconc.DoNotRequeue()
			return awsint.AWSAccountBindingApproval{}, cres, e
		}

		cres, e := reconc.RequeueWithError(err)
		return awsint.AWSAccountBindingApproval{}, cres, e
	}

	cres, e := reconc.ContinueReconciling()
	return instance, cres, e
}

// GetResources queries the API for resources necessary to determine the state
// of the existing AWSAccountBindingApproval
func (r *AWSAccountBindingApprovalReconciler) GetResources(ctx context.Context) (states.AccountBindingApprovalResources, *ctrl.Result, error) {
	instance, res, err := r.GetInstance(ctx)
	if reconc.ShouldHaltOrRequeue(res, err) {
		r.Log.Info("GetResources() halting while calling GetInstance")
		return states.AccountBindingApprovalResources{}, res, err
	}

	r.Log.Info("GetResources() completed successfully")
	cres, e := reconc.ContinueReconciling()
	return states.NewAccountBindingApprovalResources(instance), cres, e
}

// DetermineState queries the API for resources necessary to determine the state
// of existing resources, and then returns the state.
func (r *AWSAccountBindingApprovalReconciler) DetermineState(ctx context.Context) (states.AccountBindingApprovalState, *ctrl.Result, error) {
	resource, res, err := r.GetResources(ctx)
	if reconc.ShouldHaltOrRequeue(res, err) {
		r.Log.Info("DetermineState() halting while calling GetInstance")
		return states.AccountBindingApprovalState{}, res, err
	}

	r.Log.Info("DetermineState() completed successfully")
	cres, e := reconc.ContinueReconciling()
	return resource.ParseState(), cres, e
}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSAccountBindingApprovalReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&awsint.AWSAccountBindingApproval{}).
		Complete(r)
}
