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

func (r *AWSAccountBindingRequestReconciler) createAccountBindingApproval(ctx context.Context) (*ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "createAccountBindingApproval()")
	lgr.Info("starting")
	defer lgr.Info("ending")

	instanceKey := ctx.Value(instKeyContextKey).(types.NamespacedName)

	st, result, err := r.DetermineState(ctx)
	if reconc.ShouldHaltOrRequeue(result, err) {
		return result, err
	}

	// preconfigure the spec, but we'll get the resource before we create
	// and if it exists, those values will be overwritten anyway.
	approval := awsint.AWSAccountBindingApproval{
		ObjectMeta: metav1.ObjectMeta{Name: approvalInstanceKey(instanceKey).Name},
		Spec: awsint.AWSAccountBindingApprovalSpec{
			AccountID: st.AccountID(),
			ARN:       st.ARN(),
		},
	}

	if err := r.Get(ctx, approvalInstanceKey(instanceKey), &approval); err != nil {
		if apierrors.IsNotFound(err) {
			lgr.Info("creating account binding approval", "targetNamespace", instanceKey.Namespace)
			if err := r.Create(ctx, &approval); err != nil {
				lgr.Error(err, "binding creation error", "resource", instanceKey.Name)
				return reconc.RequeueWithError(err)
			}

			return reconc.ContinueReconciling()
		}

		// we got some other error
		lgr.Error(err, "some error to poll for account binding approval")
		return reconc.RequeueWithError(err)
	}

	lgr.Info("not creating resource as it already exists") // TODO more verbose logging
	return reconc.ContinueReconciling()
}

func (r *AWSAccountBindingRequestReconciler) removeAccountBindingApproval(ctx context.Context) (*ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "removeAccountBindingApproval()")
	lgr.Info("starting")
	defer lgr.Info("ending")

	instanceKey := ctx.Value(instKeyContextKey).(types.NamespacedName)

	// TODO need to generate this binding from inputs.
	binding := awsint.AWSAccountBindingApproval{
		ObjectMeta: metav1.ObjectMeta{Name: approvalInstanceKey(instanceKey).Name},
	}

	lgr.Info("deleting account binding approval", "targetNamespace", instanceKey.Namespace)
	if err := r.Delete(ctx, &binding); err != nil {
		if apierrors.IsNotFound(err) {
			lgr.Info("resource not found, it was likely deleted.")
			// continue reconciliation in case there's more to do.
			return reconc.ContinueReconciling()
		}

		lgr.Error(err, "account binding approval deletion error", "resource", instanceKey.Name)
		return reconc.RequeueWithError(err)
	}

	return reconc.ContinueReconciling()
}

// addFinalizer will add the finalizer from the instance. This function matches
// the function alias subreconcilerFuncs.
func (r *AWSAccountBindingRequestReconciler) addFinalizer(ctx context.Context) (*ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "addFinalizer")
	lgr.Info("starting")
	defer lgr.Info("ending")

	// call handler
	if res, err := r.handleFinalizer(ctx, controllerutil.AddFinalizer); reconc.ShouldHaltOrRequeue(res, err) {
		if err != nil {
			lgr.Error(err, "error handling finalizer")
		}
		return res, err
	}

	return reconc.ContinueReconciling()
}

// removeFinalizer will remove the finalizer from the instance. This function matches
// the function alias subreconcilerFuncs.
func (r *AWSAccountBindingRequestReconciler) removeFinalizer(ctx context.Context) (*ctrl.Result, error) {
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

// updateStatus manages status reconciliations.
func (r *AWSAccountBindingRequestReconciler) updateStatus(ctx context.Context) (*ctrl.Result, error) {
	lgr := r.Log.WithValues("task", "updatedStatus")
	lgr.Info("starting")
	defer lgr.Info("ending")

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
