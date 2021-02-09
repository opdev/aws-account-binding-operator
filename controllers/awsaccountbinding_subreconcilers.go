package controllers

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	ctrl "sigs.k8s.io/controller-runtime"

	reconc "github.com/opdev/aws-account-binding-operator/helpers/reconcileresults"
)

// removeFinalizer will remove the finalizer from the instance. This function matches
// the function alias subreconcilerFuncs.
func (r *AWSAccountBindingReconciler) removeFinalizer(ctx context.Context) (*ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "removeFinalizer")
	lgr.Info("starting")
	defer lgr.Info("ending")

	// call handler
	if res, err := r.handleFinalizer(ctx, controllerutil.RemoveFinalizer); reconc.ShouldHaltOrRequeue(res, err) {
		if err != nil {
			lgr.Error(err, "error handling finalizer")
		}
		return res, err
	}

	return reconc.ContinueReconciling()
}

// addFinalizer will add the finalizer from the instance. This function matches
// the function alias subreconcilerFuncs.
func (r *AWSAccountBindingReconciler) addFinalizer(ctx context.Context) (*ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "addFinalizer")
	lgr.Info("starting")
	defer lgr.Info("ending")

	// call handler
	if res, err := r.handleFinalizer(ctx, controllerutil.AddFinalizer); reconc.ShouldRequeue(res, err) {
		if err != nil {
			lgr.Error(err, "error handling finalizer")
		}
		return res, err
	}

	return reconc.ContinueReconciling()
}

// updateStatus manages status reconciliations.
func (r *AWSAccountBindingReconciler) updateStatus(ctx context.Context) (*ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "updateStatus")
	lgr.Info("starting reconciliation")
	defer lgr.Info("ending reconciliation")

	st, result, err := r.DetermineState(ctx)
	if reconc.ShouldHaltOrRequeue(result, err) {
		return result, err
	}

	inst := st.Instance()
	crStatus := inst.Status
	status := st.CurrentStatus()
	if status != crStatus {
		lgr.Info(fmt.Sprintf("updating status for resources %s", inst.GetName()))
		updated := inst.DeepCopy()
		updated.Status = status
		if err := r.Status().Update(ctx, updated); err != nil {
			lgr.Error(err, "instance patching error", "resource", inst.GetName())
			return reconc.RequeueWithError(err)
		}
	}
	return reconc.ContinueReconciling()
}

func (r *AWSAccountBindingReconciler) addNamespaceAnnotation(ctx context.Context) (*ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "addNamespaceAnnotation")
	lgr.Info("starting reconciliation")
	defer lgr.Info("ending reconciliation")

	if res, err := r.handleNamespace(ctx, addAnnotation); reconc.ShouldHaltOrRequeue(res, err) {
		return res, err
	}

	return reconc.ContinueReconciling()
}

func (r *AWSAccountBindingReconciler) removeNamespaceAnnotation(ctx context.Context) (*ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "removeNamespaceAnnotation")
	lgr.Info("starting reconciliation")
	defer lgr.Info("ending reconciliation")

	if res, err := r.handleNamespace(ctx, deleteAnnotation); reconc.ShouldHaltOrRequeue(res, err) {
		return res, err
	}

	return reconc.ContinueReconciling()
}
