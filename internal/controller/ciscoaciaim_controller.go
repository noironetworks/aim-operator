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

package controller

import (
	"context"

	ciscoaciaimv1 "github.com/noironetworks/aci-integration-module-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// GetClient -
func (r *CiscoAciAimReconciler) GetClient() client.Client {
	return r.Client
}

// GetKClient -
func (r *CiscoAciAimReconciler) GetKClient() kubernetes.Interface {
	return r.Kclient
}

// GetScheme -
func (r *CiscoAciAimReconciler) GetScheme() *runtime.Scheme {
	return r.Scheme
}

// CiscoAciAimReconciler reconciles a CiscoAciAim object
type CiscoAciAimReconciler struct {
	client.Client
	Kclient kubernetes.Interface
	Scheme  *runtime.Scheme
}

// +kubebuilder:rbac:groups=api.cisco.com,resources=ciscoaciaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=api.cisco.com,resources=ciscoaciaims/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=api.cisco.com,resources=ciscoaciaims/finalizers,verbs=update
// +kubebuilder:rbac:groups=rabbitmq.openstack.org,resources=transporturls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=topology.openstack.org,resources=topologies,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete;
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete;
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete;
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;create;update;patch;delete;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CiscoAciAim object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *CiscoAciAimReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	Log := r.GetLogger(ctx)
	instance := &ciscoaciaimv1.CiscoAciAim{}

	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if k8s_errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected.
			// For additional cleanup logic use finalizers. Return and don't requeue.
			Log.Info("CiscoAciAim instance not found, probably deleted before reconciled. Nothing to do.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		Log.Error(err, "Failed to read the Cisco Aci Aim instance.")
		return ctrl.Result{}, err
	}
	helper, err := helper.NewHelper(
		instance,
		r.Client,
		r.Kclient,
		r.Scheme,
		Log,
	)
	if err != nil {
		return ctrl.Result{}, err
	}

	//
	// initialize status
	//
	isNewInstance := instance.Status.Conditions == nil
	if isNewInstance {
		instance.Status.Conditions = condition.Conditions{}
	}

	// Save a copy of the condtions so that we can restore the LastTransitionTime
	// when a condition's state doesn't change.
	savedConditions := instance.Status.Conditions.DeepCopy()

	// Always patch the instance status when exiting this function so we can persist any changes.
	defer func() {
		// update the Ready condition based on the sub conditions
		if instance.Status.Conditions.AllSubConditionIsTrue() {
			instance.Status.Conditions.MarkTrue(
				condition.ReadyCondition, condition.ReadyMessage)
		} else {
			// something is not ready so reset the Ready condition
			instance.Status.Conditions.MarkUnknown(
				condition.ReadyCondition, condition.InitReason, condition.ReadyInitMessage)
			// and recalculate it based on the state of the rest of the conditions
			instance.Status.Conditions.Set(
				instance.Status.Conditions.Mirror(condition.ReadyCondition))
		}
		condition.RestoreLastTransitionTimes(&instance.Status.Conditions, savedConditions)
		err := h.PatchInstance(ctx, instance)
		if err != nil {
			_err = err
			return
		}
	}()
	//
	// Conditions init
	//
	cl := condition.CreateList(
		condition.UnknownCondition(condition.DBReadyCondition, condition.InitReason, condition.DBReadyInitMessage),
		condition.UnknownCondition(condition.DBSyncReadyCondition, condition.InitReason, condition.DBSyncReadyInitMessage),
		condition.UnknownCondition(condition.RabbitMqTransportURLReadyCondition, condition.InitReason, condition.RabbitMqTransportURLReadyInitMessage),
		condition.UnknownCondition(condition.CreateServiceReadyCondition, condition.InitReason, condition.CreateServiceReadyInitMessage),
		condition.UnknownCondition(condition.BootstrapReadyCondition, condition.InitReason, condition.BootstrapReadyInitMessage),
		condition.UnknownCondition(condition.InputReadyCondition, condition.InitReason, condition.InputReadyInitMessage),
		condition.UnknownCondition(condition.ServiceConfigReadyCondition, condition.InitReason, condition.ServiceConfigReadyInitMessage),
		condition.UnknownCondition(condition.DeploymentReadyCondition, condition.InitReason, condition.DeploymentReadyInitMessage),
		// service account, role, rolebinding conditions
		condition.UnknownCondition(condition.ServiceAccountReadyCondition, condition.InitReason, condition.ServiceAccountReadyInitMessage),
		condition.UnknownCondition(condition.RoleReadyCondition, condition.InitReason, condition.RoleReadyInitMessage),
		condition.UnknownCondition(condition.RoleBindingReadyCondition, condition.InitReason, condition.RoleBindingReadyInitMessage),
	)

	instance.Status.Conditions.Init(&cl)
	instance.Status.ObservedGeneration = instance.Generation

	// If we're not deleting this and the service object doesn't have our finalizer, add it.
	if instance.DeletionTimestamp.IsZero() && controllerutil.AddFinalizer(instance, helper.GetFinalizer()) || isNewInstance {
		return ctrl.Result{}, nil
	}

	if instance.Status.Hash == nil {
		instance.Status.Hash = map[string]string{}
	}

	// Init Topology condition if there's a reference
	if instance.Spec.TopologyRef != nil {
		c := condition.UnknownCondition(condition.TopologyReadyCondition, condition.InitReason, condition.TopologyReadyInitMessage)
		cl.Set(c)
	}

	// Handle service delete
	if !instance.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, instance, helper)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, instance, helper)
}

// fields to index to reconcile when change
const (
	passwordSecretField      = ".spec.secret"
	rabbitmqClusterNameField = ".spec.rabbitmqClusterName"
	topologyField            = ".spec.topologyRef.Name"
)

// SetupWithManager sets up the controller with the Manager.
func (r *CiscoAciAimReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.CiscoAciAim{}).
		Complete(r)
}

func (r *CiscoAciAimReconciler) generateServiceConfig(
	ctx context.Context,
	h *helper.Helper,
	instance *ciscoaciaimv1.CiscoAciAim,
	envVars *map[string]env.Setter,
	serviceLabels map[string]string,
	db *mariadbv1.Database,
) error {
	Log := r.GetLogger(ctx)
	Log.Info("generateServiceConfigMaps - CiscoAciAim controller")

	// create Secret required for barbican input
	labels := labels.GetLabels(instance, labels.GetGroupLabel(ciscoaciaim.ServiceName), serviceLabels)

	customData := map[string]string{
		ciscoaciaim.CustomConfigFileName: instance.Spec.CustomServiceConfig,
	}

	for key, data := range instance.Spec.DefaultConfigOverwrite {
		customData[key] = data
	}

	cms := []util.Template{
		// ScriptsConfigMap
		{
			Name:         fmt.Sprintf("%s-scripts", instance.Name),
			Namespace:    instance.Namespace,
			Type:         util.TemplateTypeScripts,
			InstanceType: instance.Kind,
			Labels:       cmLabels,
		},
		// ConfigMap
		{
			Name:               fmt.Sprintf("%s-config-data", instance.Name),
			Namespace:          instance.Namespace,
			Type:               util.TemplateTypeConfig,
			InstanceType:       instance.Kind,
			CustomData:         customData,
			ConfigOptions:      templateParameters,
			Labels:             cmLabels,
			AdditionalTemplate: extraTemplates,
		},
	}
	return secret.EnsureSecrets(ctx, h, instance, cms, envVars)
}
