/*
Copyright 2026.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	computev1 "github.com/shkatara/operator/api/v1"
)

// EC2instanceReconciler reconciles a EC2instance object
type EC2instanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=compute.cloud.com,resources=ec2instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=compute.cloud.com,resources=ec2instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=compute.cloud.com,resources=ec2instances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EC2instance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.24.1/pkg/reconcile
func (r *EC2instanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := logf.FromContext(ctx)
	ec2Instance := &computev1.EC2instance{}
	r.Get(ctx, req.NamespacedName, ec2Instance)
	fmt.Println("Reconcile called for EC2instance:", ec2Instance.Spec.InstanceName)
	fmt.Println("Spec:", ec2Instance.Spec)
	fmt.Println("Status:", ec2Instance.Status)
	fmt.Println("InstanceType:", ec2Instance.Status.InstanceType)
	fmt.Println("Tags:", ec2Instance.Spec.Tags)
	fmt.Println("Storage:", ec2Instance.Spec.Storage)
	fmt.Println("SshKey:", ec2Instance.Spec.SshKey)
	fmt.Println("AMI ID:", ec2Instance.Spec.AmiId)
	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EC2instanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&computev1.EC2instance{}).
		Named("ec2instance").
		Complete(r)
}
