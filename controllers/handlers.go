package controllers

import (
	"github.com/opdev/aws-account-binding-operator/controllers/constants"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// finalizerFunc abstracts the function signature of controllerutil.AddFinalizer
// and controllerutil.RemoveFinalizer
type finalizerFunc func(o client.Object, finalizer string)

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
