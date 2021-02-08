package controllers

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/opdev/aws-account-binding-operator/helpers/reconcileresults"
)

// removeFinalizer will remove the finalizer from the instance. This function matches
// the function alias subreconcilerFuncs.
func (r *AWSAccountBindingReconciler) removeFinalizer(ctx context.Context) (*ctrl.Result, error) {
	logger := r.Log.WithName("removeFinalizer")
	logger.Info("starting")
	defer logger.Info("ending")

	// call handler
	if res, err := r.handleFinalizer(ctx, controllerutil.RemoveFinalizer); ShouldHaltOrRequeue(res, err) {
		if err != nil {
			logger.Error(err, "error handling finalizer")
		}
		return res, err
	}

	return ContinueReconciling()
}

// addFinalizer will add the finalizer from the instance. This function matches
// the function alias subreconcilerFuncs.
func (r *AWSAccountBindingReconciler) addFinalizer(ctx context.Context) (*ctrl.Result, error) {
	logger := r.Log.WithName("addFinalizer")
	logger.Info("starting")
	defer logger.Info("ending")

	// call handler
	if res, err := r.handleFinalizer(ctx, controllerutil.AddFinalizer); ShouldRequeue(res, err) {
		if err != nil {
			logger.Error(err, "error handling finalizer")
		}
		return res, err
	}

	return ContinueReconciling()
}

// updateStatus manages status reconciliations.
func (r *AWSAccountBindingReconciler) updateStatus(ctx context.Context) (*ctrl.Result, error) {
	logger := r.Log.WithName("updateStatus")
	logger.Info("starting reconciliation")
	defer logger.Info("ending reconciliation")

	st, result, err := r.DetermineState(ctx)
	if ShouldHaltOrRequeue(result, err) {
		return result, err
	}

	inst := st.Instance()
	crStatus := inst.Status
	status := st.CurrentStatus()
	if status != crStatus {
		logger.Info(fmt.Sprintf("updating status for resources %s", inst.GetName()))
		updated := inst.DeepCopy()
		updated.Status = status
		if err := r.Status().Update(ctx, updated); err != nil {
			logger.Error(err, "instance patching error", "resource", inst.GetName())
			return RequeueWithError(err)
		}
	}
	return ContinueReconciling()
}

func (r *AWSAccountBindingReconciler) addNamespaceAnnotation(ctx context.Context) (*ctrl.Result, error) {
	logger := r.Log.WithName("addNamespaceAnnotation")
	logger.Info("starting reconciliation")
	defer logger.Info("ending reconciliation")

	if res, err := r.handleNamespace(ctx, addAnnotation); ShouldHaltOrRequeue(res, err) {
		return res, err
	}

	return ContinueReconciling()
}

func (r *AWSAccountBindingReconciler) removeNamespaceAnnotation(ctx context.Context) (*ctrl.Result, error) {
	logger := r.Log.WithName("removeNamespaceAnnotation")
	logger.Info("starting reconciliation")
	defer logger.Info("ending reconciliation")

	if res, err := r.handleNamespace(ctx, deleteAnnotation); ShouldHaltOrRequeue(res, err) {
		return res, err
	}

	return ContinueReconciling()
}
