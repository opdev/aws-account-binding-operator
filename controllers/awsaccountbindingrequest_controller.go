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

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	awsint "github.com/opdev/aws-account-binding-operator/api/v1alpha1"
	"github.com/opdev/aws-account-binding-operator/controllers/states"
	reconc "github.com/opdev/aws-account-binding-operator/helpers/reconcileresults"
)

// AWSAccountBindingRequestReconciler reconciles a AWSAccountBindingRequest object
type AWSAccountBindingRequestReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=integrations.opdev.io,resources=awsaccountbindingrequests,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=integrations.opdev.io,resources=awsaccountbindingrequests/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=integrations.opdev.io,resources=awsaccountbindingrequests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *AWSAccountBindingRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = context.WithValue(ctx, instKeyContextKey, req.NamespacedName)
	lgr := r.Log.WithValues("Request", req.NamespacedName)
	lgr.Info("starting reconciliation")
	defer lgr.Info("ending reconciliation")

	state, res, err := r.DetermineState(ctx)
	if reconc.ShouldHaltOrRequeue(res, err) {
		lgr.Info("Reconcile() halting while calling DetermineState")
		return reconc.Evaluate(res, err)
	}

	// handle deletion
	deletionSubReconcilers := []subreconcilerFuncs{
		r.removeAccountBindingApproval,
		r.removeFinalizer,
	}

	if state.IsBeingDeleted() {
		lgr.Info("resource is being deleted, running deletion reconciliation flows")
		for _, f := range deletionSubReconcilers {
			if r, err := f(ctx); reconc.ShouldHaltOrRequeue(r, err) {
				return reconc.Evaluate(r, err)
			}
		}

		return reconc.Evaluate(reconc.DoNotRequeue())
	}

	lgr.Info("running reconciliation flows")
	subreconcilers := []subreconcilerFuncs{
		r.updateStatus,
		r.addFinalizer,
		r.createAccountBindingApproval,
	}

	for _, f := range subreconcilers {
		// call the reconciler with the state
		if r, err := f(ctx); reconc.ShouldHaltOrRequeue(r, err) {
			return reconc.Evaluate(r, err)
		}
	}

	return ctrl.Result{}, nil
}

// GetInstance queries the API for the instance of the custom resource.
func (r *AWSAccountBindingRequestReconciler) GetInstance(ctx context.Context) (awsint.AWSAccountBindingRequest, *ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "GetInstance")
	lgr.Info("starting")
	defer lgr.Info("ending")

	instanceKey := ctx.Value(instKeyContextKey).(types.NamespacedName)
	var instance awsint.AWSAccountBindingRequest
	if err := r.Get(ctx, instanceKey, &instance); err != nil {
		if apierrors.IsNotFound(err) {
			// it was deleted before reconcile completed
			lgr.Info("GetInstance() resource not found, it was likely deleted.")
			cres, e := reconc.DoNotRequeue()
			return awsint.AWSAccountBindingRequest{}, cres, e
		}

		cres, e := reconc.RequeueWithError(err)
		return awsint.AWSAccountBindingRequest{}, cres, e
	}

	cres, e := reconc.ContinueReconciling()
	return instance, cres, e
}

// GetApprovalInstance queries the API for the instance of the AWSAccountBindingApproval
// associated with a given AWSAccountBindingRequest
func (r *AWSAccountBindingRequestReconciler) GetApprovalInstance(ctx context.Context) (awsint.AWSAccountBindingApproval, *ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "GetApprovalIntance")
	lgr.Info("starting")
	defer lgr.Info("ending")

	approvalKey := approvalInstanceKey(ctx.Value(instKeyContextKey).(types.NamespacedName))
	var instance awsint.AWSAccountBindingApproval
	if err := r.Get(ctx, approvalKey, &instance); err != nil {
		if apierrors.IsNotFound(err) {
			// it was deleted before reconcile completed
			lgr.Info("GetApprovalInstance() resource not found, it was likely deleted.")
			// we are okay with it being deleted before we complete, so we don't
			// requeue here. If it doesn't exist, then the approval status
			// is false.
			cres, e := reconc.ContinueReconciling()
			return awsint.AWSAccountBindingApproval{}, cres, e
		}

		// There was a problem with getting the instance
		// we might not need to requeue here either, but
		// for now we'll requeue and try again.
		cres, e := reconc.RequeueWithError(err)
		return awsint.AWSAccountBindingApproval{}, cres, e
	}

	cres, e := reconc.ContinueReconciling()
	return instance, cres, e
}

// GetResources queries the API for resources necessary to determine the state
// of the existing AWSAccountBindingRequest
func (r *AWSAccountBindingRequestReconciler) GetResources(ctx context.Context) (states.AccountBindingRequestResources, *ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "GetResources")
	lgr.Info("starting")
	defer lgr.Info("ending")

	instance, res, err := r.GetInstance(ctx)
	if reconc.ShouldHaltOrRequeue(res, err) {
		lgr.Info("GetResources() halting while calling GetInstance")
		return states.AccountBindingRequestResources{}, res, err
	}

	approvalInstance, res, err := r.GetApprovalInstance(ctx)
	if reconc.ShouldHaltOrRequeue(res, err) {
		lgr.Info("GetResources() halting while calling GetApprovalInstance")
		return states.AccountBindingRequestResources{}, res, err
	}

	lgr.Info("GetResources() completed successfully")
	cres, e := reconc.ContinueReconciling()
	return states.NewAccountBindingRequestResources(instance, approvalInstance.Status), cres, e
}

// DetermineState queries the API for resources necessary to determine the state
// of existing resources, and then returns the state.
func (r *AWSAccountBindingRequestReconciler) DetermineState(ctx context.Context) (states.AccountBindingRequestState, *ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "DetermineState")
	lgr.Info("starting")
	defer lgr.Info("ending")

	resource, res, err := r.GetResources(ctx)
	if reconc.ShouldHaltOrRequeue(res, err) {
		lgr.Info("DetermineState() halting while calling GetInstance")
		return states.AccountBindingRequestState{}, res, err
	}

	lgr.Info("DetermineState() completed successfully")
	cres, e := reconc.ContinueReconciling()
	return resource.ParseState(), cres, e
}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSAccountBindingRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&awsint.AWSAccountBindingRequest{}).
		Complete(r)
}
