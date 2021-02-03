package controllers

import (
	"context"

	"github.com/opdev/aws-account-binding-operator/controllers/constants"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/opdev/aws-account-binding-operator/helpers/reconcileresults"
)

// finalizerFunc abstracts the function signature of controllerutil.AddFinalizer
// and controllerutil.RemoveFinalizer
type finalizerFunc func(o client.Object, finalizer string)

// handleFinalizer executes finalizer management
func (r *AWSAccountBindingReconciler) handleFinalizer(ctx context.Context, modifyFinalizer finalizerFunc) (*ctrl.Result, error) {
	resources, result, err := r.GetResources(ctx)
	if ShouldHaltOrRequeue(result, err) {
		return result, err
	}

	inst := resources.Instance()
	patch := client.MergeFrom(inst.DeepCopy())
	modifyFinalizer(inst, constants.Finalizer)
	if err := r.Patch(ctx, inst, patch); err != nil {
		return RequeueWithError(err)
	}

	return ContinueReconciling()
}

type annotationFunc func(map[string]string, string) map[string]string

func deleteAnnotation(annot map[string]string, val string) map[string]string {
	// while this accepts val, it does not use it for the
	// delete operation. The function signature is kept
	// consistent with the addAnnotation function.
	delete(annot, constants.Annotation)
	return annot
}

func addAnnotation(annot map[string]string, val string) map[string]string {
	if annot == nil {
		// if the input annotation map is nil
		annot = make(map[string]string)
	}

	annot[constants.Annotation] = val
	return annot
}

// handleNamespace handles
func (r *AWSAccountBindingReconciler) handleNamespace(ctx context.Context, handleAnnotation annotationFunc) (*ctrl.Result, error) {
	st, result, err := r.DetermineState(ctx)
	if ShouldHaltOrRequeue(result, err) {
		return result, err
	}

	ns := st.Namespace()
	inst := st.Instance()
	annotationHandlerVal := inst.Spec.AccountID

	baseToPatch := client.MergeFrom(ns.DeepCopy())
	ns.Annotations = handleAnnotation(ns.GetAnnotations(), annotationHandlerVal)
	if err := r.Patch(ctx, ns, baseToPatch); err != nil {
		return RequeueWithError(err)
	}

	return ContinueReconciling()
}
