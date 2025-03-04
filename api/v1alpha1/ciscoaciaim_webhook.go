/*
Copyright 2025.

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
package ciscoaciaim


import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
    apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	webhook "sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
    topologyv1 "github.com/openstack-k8s-operators/infra-operator/apis/topology/v1beta1"
)

// log is for logging in this package.
var ciscoaciaimlog = logf.Log.WithName("ciscoaciaim-resource")

// CiscoAciAimDefaults -
type CiscoAciAimDefaults struct {
	ContainerImageURL string
	APITimeout        int
}

var ciscoAciAimDefaults CiscoAciAimDefaults


// SetupCiscoAciAimDefaults - initialize CiscoAciAimAPI spec defaults for use with either internal or external webhooks
func SetupCiscoAciAimDefaults(defaults CiscoAciAimDefaults) {
	ciscoAciAimDefaults = defaults

	ciscoaciaimlog.Info("CiscoAciAimDefaults defaults initialized", "defaults", defaults)
}


// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *CiscoAciAim) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-api-cisco-com-v1alpha1-ciscoaciaim,mutating=true,failurePolicy=fail,sideEffects=None,groups=api.cisco.com,resources=ciscoaciaims,verbs=create;update,versions=v1alpha1,name=mciscoaciaim.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &CiscoAciAim{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *CiscoAciAim) Default() {
	ciscoaciaimlog.Info("default", "name", r.Name)


	if r.Spec.ContainerImage == "" {
		r.Spec.ContainerImage = ciscoAciAimDefaults.ContainerImageURL
	}
	r.Spec.Default()

}

// Default - set defaults for this KeystoneAPI spec
func (spec *CiscoAciAimSpec) Default() {
	// no defaults to set yet
	spec.CiscoAciAimSpecCore.Default()
}

func (spec *CiscoAciAimSpecCore) Default() {
	// nothing here yet
}


// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-api-cisco-com-v1alpha1-ciscoaciaim,mutating=false,failurePolicy=fail,sideEffects=None,groups=api.cisco.com,resources=ciscoaciaims,verbs=create;update,versions=v1alpha1,name=vciscoaciaim.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &CiscoAciAim{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *CiscoAciAim) ValidateCreate() (admission.Warnings, error) {
	ciscoaciaimlog.Info("validate create", "name", r.Name)

	var allErrs field.ErrorList
	basePath := field.NewPath("spec")

	// When a TopologyRef CR is referenced, fail if a different Namespace is
	// referenced because is not supported
	if r.Spec.TopologyRef != nil {
		if err := topologyv1.ValidateTopologyNamespace(r.Spec.TopologyRef.Namespace, *basePath, r.Namespace); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	if len(allErrs) != 0 {
		return nil, apierrors.NewInvalid(
            GroupVersion.WithKind("CiscoAciAim").GroupKind(), 
			r.Name, allErrs)
	}
	return nil, nil

	// TODO(user): fill in your validation logic upon object creation.
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *CiscoAciAim) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	ciscoaciaimlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *CiscoAciAim) ValidateDelete() (admission.Warnings, error) {
	ciscoaciaimlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
