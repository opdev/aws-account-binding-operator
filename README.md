# AWS Account Binding Operator

The AWS Account Binding Operator is a series of controllers managing the binding
of an AWS Account to a Namespace in an OpenShift or Kubernetes cluster. This binding
is designed to function in conjunction with the [AWS Controllers for Kubernetes](https://github.com/aws-controllers-k8s)
("ACK").

This project is a **proof of concept** and is not directly related to the ACK project.
Use of this project is done at your own risk.

## The Motive

The [AWS Controllers for Kubernetes](https://github.com/aws-controllers-k8s) are
designed to enable mulit-tenancy through the use of
[Cross-Account Resource Management](https://github.com/aws-controllers-k8s/community/blob/main/docs/design/proposals/carm/cross-account-resource-management.md)
("CARM").

The implementation for the AWS service controllers to effectively pivot API calls
to be issued against specified account. This is done via the `services.k8s.aws/owner-account-id:`
annotation, which lives on a given Namespace resource.

When a cluster administrator has applied this annotation to a namespace, the AWS
service controller will read a configmap and search for a corresponding accountID:ARN
mapping, which informs it how to pivot, allowing it to issue API calls on behalf
of the requested account.

## The Concept

This operator allows for any non cluster-administrator to request for an AWS
account to be bound to their namespace. A cluster-administrator can then acknowledge
and approve the request. Once a request is approved, the operator will annotate
the given namespace with the desired AWS account.

When a user deletes the request, the associated approval and and binding are both
deleted as well.

## The Implementation

The operator introduces three custom resource definitions.

An **AWSAccountBindingRequest** represents a user's request to have a given AWS account
associated with their namespace. This is a namespace-scoped resource. This resource
should be fully manageable by regular users within the associated namespace.

An **AWSAccountBindingApproval** represents an administrators explicit approval
of a given **AWSAccountBindingRequest**. This is a cluster-scoped resource which
takes the name of the resource requesting the binding. This ensures that a single
instance of a given AWSAccountBindingApproval can exist for any namespace. This
resource is created by a controller, and is expected to be edited by a human.

An **AWSAccountBinding** represents a managed annotation on a given namespace.
When an **AWSAccountBindingApproval** is approved, a controller will automatically
create this resource, which then triggers the annotation of the namespace.

NOTE: At the moment, there is no management of the associated configmap that the
ACK expects. This is a required feature expected in the short-term future.