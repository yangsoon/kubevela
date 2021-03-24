/*
 Copyright 2021. The KubeVela Authors.

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

package v1beta1

import (
	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/common"
	"github.com/oam-dev/kubevela/apis/standard.oam.dev/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AppStatus defines the observed state of Application
type AppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	v1alpha1.RolloutStatus `json:",inline"`

	Phase common.ApplicationPhase `json:"status,omitempty"`

	// Components record the related Components created by Application Controller
	Components []runtimev1alpha1.TypedReference `json:"components,omitempty"`

	// Services record the status of the application services
	Services []common.ApplicationComponentStatus `json:"services,omitempty"`

	// LatestRevision of the application configuration it generates
	// +optional
	LatestRevision *common.Revision `json:"latestRevision,omitempty"`
}

// ApplicationComponent describe the component of application
type ApplicationComponent struct {
	Name         string `json:"name"`
	WorkloadType string `json:"type"`
	// +kubebuilder:pruning:PreserveUnknownFields
	Properties runtime.RawExtension `json:"properties"`

	// Traits define the trait of one component, the type must be array to keep the order.
	Traits []common.ApplicationTrait `json:"traits,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// scopes in ApplicationComponent defines the component-level scopes
	// the format is <scope-type:scope-instance-name> pairs, the key represents type of `ScopeDefinition` while the value represent the name of scope instance.
	Scopes map[string]string `json:"scopes,omitempty"`
}

// ApplicationSpec is the spec of Application
type ApplicationSpec struct {
	Components []ApplicationComponent `json:"components"`

	// TODO(wonderflow): we should have application level scopes supported here

	// RolloutPlan is the details on how to rollout the resources
	// The controller simply replace the old resources with the new one if there is no rollout plan involved
	// +optional
	RolloutPlan *v1alpha1.RolloutPlan `json:"rolloutPlan,omitempty"`
}

// +kubebuilder:object:root=true

// Application is the Schema for the applications API
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec `json:"spec,omitempty"`
	Status AppStatus       `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

// GetComponent get the component from the application based on its workload type
func (app *Application) GetComponent(workloadType string) *ApplicationComponent {
	for _, c := range app.Spec.Components {
		if c.WorkloadType == workloadType {
			return &c
		}
	}
	return nil
}
