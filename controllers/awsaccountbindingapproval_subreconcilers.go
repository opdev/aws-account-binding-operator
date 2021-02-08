package controllers

import (
	"context"
	"fmt"

	awsint "github.com/opdev/aws-account-binding-operator/api/v1alpha1"
	reconc "github.com/opdev/aws-account-binding-operator/helpers/reconcileresults"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *AWSAccountBindingApprovalReconciler) ensureApproval(ctx context.Context) (*ctrl.Result, error) {
	logger := r.Log.WithName("haltIfNotApproved()")
	logger.Info("starting reconciliation")
	defer logger.Info("ending reconciliation")

	state, res, err := r.DetermineState(ctx)
	if reconc.ShouldHaltOrRequeue(res, err) {
		r.Log.Info("halting while calling DetermineState")
		return reconc.DoNotRequeue()
	}

	if !state.IsApproved() {
		instanceKey := ctx.Value(instKeyContextKey).(types.NamespacedName)
		r.Log.Info("halting because the request is not approved", "resource", instanceKey)
		return reconc.DoNotRequeue()
	}

	return reconc.ContinueReconciling()
}

// removeFinalizer will remove the finalizer from the instance. This function matches
// the function alias subreconcilerFuncs.
// TODO lots of duplciate code here with other controllers in this operator.
// restructure these to register generic functions as individual struct methods.
func (r *AWSAccountBindingApprovalReconciler) removeFinalizer(ctx context.Context) (*ctrl.Result, error) {
	logger := r.Log.WithName("removeFinalizer")
	logger.Info("starting")
	defer logger.Info("ending")

	// call handler
	if res, err := r.handleFinalizer(ctx, controllerutil.RemoveFinalizer); reconc.ShouldHaltOrRequeue(res, err) {
		if err != nil {
			logger.Error(err, "error handling finalizer")
		}
		return res, err
	}

	return reconc.ContinueReconciling()
}

// addFinalizer will add the finalizer from the instance. This function matches
// the function alias subreconcilerFuncs.
func (r *AWSAccountBindingApprovalReconciler) addFinalizer(ctx context.Context) (*ctrl.Result, error) {
	logger := r.Log.WithName("addFinalizer")
	logger.Info("starting")
	defer logger.Info("ending")

	// call handler
	if res, err := r.handleFinalizer(ctx, controllerutil.AddFinalizer); reconc.ShouldRequeue(res, err) {
		if err != nil {
			logger.Error(err, "error handling finalizer")
		}
		return res, err
	}

	return reconc.ContinueReconciling()
}

// updateStatus manages status reconciliations.
func (r *AWSAccountBindingApprovalReconciler) updateStatus(ctx context.Context) (*ctrl.Result, error) {
	logger := r.Log.WithName("updateStatus")
	logger.Info("starting reconciliation")
	defer logger.Info("ending reconciliation")

	st, result, err := r.DetermineState(ctx)
	if reconc.ShouldHaltOrRequeue(result, err) {
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
			return reconc.RequeueWithError(err)
		}
	}
	return reconc.ContinueReconciling()
}

func (r *AWSAccountBindingApprovalReconciler) createAccountBinding(ctx context.Context) (*ctrl.Result, error) {
	logger := r.Log.WithValues("task", "createAccountBinding()")
	logger.Info("starting")
	defer logger.Info("ending")

	instanceKey := ctx.Value(instKeyContextKey).(types.NamespacedName)

	st, result, err := r.DetermineState(ctx)
	if reconc.ShouldHaltOrRequeue(result, err) {
		return result, err
	}

	// preconfigure the spec, but we'll get the resource before we create
	// and if it exists, those values will be overwritten anyway.
	binding := awsint.AWSAccountBinding{
		ObjectMeta: metav1.ObjectMeta{Name: parseBindingInstanceKey(instanceKey).Name},
		Spec: awsint.AWSAccountBindingSpec{
			AccountID: st.AccountID(),
			ARN:       st.ARN(),
		},
	}

	if err := r.Get(ctx, parseBindingInstanceKey(instanceKey), &binding); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("creating account binding", "targetNamespace", instanceKey.Name)
			if err := r.Create(ctx, &binding); err != nil {
				logger.Error(err, "binding creation error", "resource", instanceKey.Name)
				return reconc.RequeueWithError(err)
			}
			// we got some other error
			logger.Error(err, "createAccountBinding() some error to poll for account binding")
			return reconc.RequeueWithError(err)
		}

		logger.Info("createAccountBinding() not creating resource as it already exists") // TODO more verbose logging
	}

	return reconc.ContinueReconciling()
}

func (r *AWSAccountBindingApprovalReconciler) removeAccountBinding(ctx context.Context) (*ctrl.Result, error) {
	logger := r.Log.WithValues("task", "removeAccountBinding()")
	logger.Info("starting")
	defer logger.Info("ending")

	instanceKey := ctx.Value(instKeyContextKey).(types.NamespacedName)

	// TODO need to generate this binding from inputs.
	binding := awsint.AWSAccountBinding{
		ObjectMeta: metav1.ObjectMeta{Name: parseBindingInstanceKey(instanceKey).Name},
	}

	logger.Info("deleting account binding", "targetNamespace", instanceKey.Name)
	if err := r.Delete(ctx, &binding); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("removeAccountBinding() resource not found, it was likely deleted.")
			// continue reconciliation in case there's more to do.
			return reconc.ContinueReconciling()
		}

		logger.Error(err, "binding deletion error", "resource", instanceKey.Name)
		return reconc.RequeueWithError(err)
	}

	return reconc.ContinueReconciling()
}
