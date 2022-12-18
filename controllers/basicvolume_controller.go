/*
Copyright 2023.

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
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	basicv1 "github.com/ragul28/kube-operator-basic/api/v1"
)

// BasicVolumeReconciler reconciles a BasicVolume object
type BasicVolumeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=basic.sample-operator.dev,resources=basicvolumes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=basic.sample-operator.dev,resources=basicvolumes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=basic.sample-operator.dev,resources=basicvolumes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BasicVolume object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *BasicVolumeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("Enter Reconcile", "req", req)

	volume := &basicv1.BasicVolume{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, volume)

	l.Info("Enter Reconcile", "spec", volume.Spec, "status", volume.Status)

	if volume.Spec.Name != volume.Status.Name {
		volume.Status.Name = volume.Spec.Name
		r.Status().Update(ctx, volume)
	}

	return ctrl.Result{}, nil
}

func (r *BasicVolumeReconciler) reconcilePVC(ctx context.Context, volume *basicv1.BasicVolume, l logr.Logger) error {
	pvc := &v1.PersistentVolumeClaim{}
	err := r.Get(ctx, types.NamespacedName{Name: volume.Name, Namespace: volume.Namespace}, pvc)
	if err == nil {
		l.Info("PVC Found")
		return nil
	}

	if !errors.IsNotFound(err) {
		return err
	}

	l.Info("PVC Not found")

	storageClass := "default"
	storageRequest, _ := resource.ParseQuantity(fmt.Sprintf("%dGi", volume.Spec.Size))

	pvc = &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: volume.Namespace,
			Name:      volume.Name,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClass,
			AccessModes:      []v1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{"storage": storageRequest},
			},
		},
	}

	l.Info("Creating PVC")
	return r.Create(ctx, pvc)

}

// SetupWithManager sets up the controller with the Manager.
func (r *BasicVolumeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&basicv1.BasicVolume{}).
		Complete(r)
}
