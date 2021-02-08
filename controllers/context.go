package controllers

// contextKey is used to pass the instance key in reconciliation contexts
type contextKey string

// instKeyContextKey is the key used to describe the resource query key
// associated with a request passed into reconciliation. The value is expected
// to fulfill types.NamespacedName which is a part of the API CRUD
// operation function signature.
var instKeyContextKey contextKey = "requested-instance-key"
