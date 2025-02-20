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

package v1alpha1

import (
        condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
        "github.com/openstack-k8s-operators/lib-common/modules/storage"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        //topologyv1 "github.com/openstack-k8s-operators/infra-operator/apis/topology/v1beta1"
)

const (
    // ContainerImage - default fall-back container image for CiscoAciAim if associated env var not provided
    ContainerImage = "quay.io/podified-antelope-centos9/openstack-ciscoaci-aim:current-podified"
)

type CiscoAciAimSpecCore struct {
        // +kubebuilder:validation:Optional
        // NodeSelector to target subset of worker nodes running this service
        NodeSelector *map[string]string `json:"nodeSelector,omitempty"`

        // +kubebuilder:validation:Optional
        // ConfigOverwrite - interface to overwrite default config files like e.g. logging.conf or policy.json.
        // But can also be used to add additional files. Those get added to the service config dir in /etc/<service> .
        DefaultConfigOverwrite map[string]string `json:"defaultConfigOverwrite,omitempty"`

        // +kubebuilder:validation:Optional
        // +kubebuilder:default=1
        // +kubebuilder:validation:Maximum=32
        // +kubebuilder:validation:Minimum=0
        // Replicas of horizon API to run
        Replicas *int32 `json:"replicas"`

        // ExtraMounts containing conf files
        // +kubebuilder:default={}
        ExtraMounts []CiscoAciAimExtraVolMounts `json:"extraMounts,omitempty"`

        // +kubebuilder:validation:Optional
        // CustomServiceConfig - customize the service config using this parameter to change service defaults,
        // or overwrite rendered information using raw OpenStack config format. The content gets added to
        // to /etc/<service>/<service>.conf.d directory as custom.conf file.
        CustomServiceConfig string `json:"customServiceConfig,omitempty"`

        // +kubebuilder:validation:Optional
        // +kubebuilder:default=neutron
        // ServiceUser - optional username used for this service to register in neutron
        ServiceUser string `json:"serviceUser"`

        // +kubebuilder:validation:Required
        // MariaDB instance name
        // Right now required by the maridb-operator to get the credentials from the instance to create the DB
        // Might not be required in future
        DatabaseInstance string `json:"databaseInstance"`

        // +kubebuilder:validation:Optional
        // +kubebuilder:default=neutron
        // DatabaseAccount - optional MariaDBAccount CR name used for neutron DB, defaults to neutron
        DatabaseAccount string `json:"databaseAccount"`

        // +kubebuilder:validation:Required
        // +kubebuilder:default=rabbitmq
        // RabbitMQ instance name
        // Needed to request a transportURL that is created and used in Neutron
        RabbitMqClusterName string `json:"rabbitMqClusterName"`

        /*
        // +kubebuilder:validation:Optional
        // TopologyRef to apply the Topology defined by the associated CR referenced
        // by name
        TopologyRef *topologyv1.TopoRef `json:"topologyRef,omitempty"`
        */
}

// CiscoAciAimSpec defines the desired state of CiscoAciAim
type CiscoAciAimSpec struct {
        // +kubebuilder:validation:Required
        // Cisco Aci Aim Container Image URL
        ContainerImage string `json:"containerImage"`
        CiscoAciAimSpecCore `json:",inline"`
}

// CiscoAciAimStatus defines the observed state of CiscoAciAim
type CiscoAciAimStatus struct {

   // Map of hashes to track e.g. job status
   Hash map[string]string `json:"hash,omitempty"`

   // Conditions
   Conditions condition.Conditions `json:"conditions,omitempty" optional:"true"`

   // Neutron Database Hostname
   DatabaseHostname string `json:"databaseHostname,omitempty"`

   // TransportURLSecret - Secret containing RabbitMQ transportURL
   TransportURLSecret string `json:"transportURLSecret,omitempty"`

   // ObservedGeneration - the most recent generation observed for this
   // service. If the observed generation is less than the spec generation,
   // then the controller has not processed the latest changes injected by
   // the opentack-operator in the top-level CR (e.g. the ContainerImage)
   ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CiscoAciAim is the Schema for the ciscoaciaims API
type CiscoAciAim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CiscoAciAimSpec   `json:"spec,omitempty"`
	Status CiscoAciAimStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CiscoAciAimList contains a list of CiscoAciAim
type CiscoAciAimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CiscoAciAim `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CiscoAciAim{}, &CiscoAciAimList{})
}

func SetupDefaults() {
    ciscoAciAimDefaults := CiscoAciAimDefaults{
        ContainerImageURL: util.GetEnvVar("RELATED_IMAGE_CISCOACI_AIM_IMAGE_URL_DEFAULT", ContainerImage),
    }

}

// CiscoAciAimExtraVolMounts exposes additional parameters processed by the horizon-operator
// and defines the common VolMounts structure provided by the main storage module
type CiscoAciAimExtraVolMounts struct {
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	Region string `json:"region,omitempty"`
	// +kubebuilder:validation:Required
	VolMounts []storage.VolMounts `json:"extraVol"`
}

// Propagate is a function used to filter VolMounts according to the specified
// PropagationType array
func (c *CiscoAciAimExtraVolMounts) Propagate(svc []storage.PropagationType) []storage.VolMounts {
	var vl []storage.VolMounts
	for _, gv := range c.VolMounts {
		vl = append(vl, gv.Propagate(svc)...)
	}
	return vl
}
