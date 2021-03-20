/*
Copyright 2021 The KubeVela Authors.

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

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1alpha2"
)

// ComponentDefinitionSpec defines the desired state of ComponentDefinition
type ComponentDefinitionSpec struct {
	// Workload is a workload type descriptor
	Workload v1alpha2.WorkloadTypeDescriptor `json:"workload"`

	// ChildResourceKinds are the list of GVK of the child resources this workload generates
	ChildResourceKinds []v1alpha2.ChildResourceKind `json:"childResourceKinds,omitempty"`

	// RevisionLabel indicates which label for underlying resources(e.g. pods) of this workload
	// can be used by trait to create resource selectors(e.g. label selector for pods).
	// +optional
	RevisionLabel string `json:"revisionLabel,omitempty"`

	// PodSpecPath indicates where/if this workload has K8s podSpec field
	// if one workload has podSpec, trait can do lot's of assumption such as port, env, volume fields.
	// +optional
	PodSpecPath string `json:"podSpecPath,omitempty"`

	// Status defines the custom health policy and status message for workload
	// +optional
	Status *v1alpha2.Status `json:"status,omitempty"`

	// Schematic defines the data format and template of the encapsulation of the workload
	// +optional
	Schematic *v1alpha2.Schematic `json:"schematic,omitempty"`

	// Extension is used for extension needs by OAM platform builders
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Extension *runtime.RawExtension `json:"extension,omitempty"`
}

// ComponentDefinitionStatus is the status of ComponentDefinition
type ComponentDefinitionStatus struct {
	// ConditionedStatus reflects the observed status of a resource
	runtimev1alpha1.ConditionedStatus `json:",inline"`
	// ConfigMapRef refer to a ConfigMap which contains OpenAPI V3 JSON schema of Component parameters.
	ConfigMapRef string `json:"configMapRef,omitempty"`
}

// +kubebuilder:object:root=true

// ComponentDefinition is the Schema for the componentdefinitions API
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,oam}
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
type ComponentDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentDefinitionSpec   `json:"spec,omitempty"`
	Status ComponentDefinitionStatus `json:"status,omitempty"`
}

// SetConditions set condition for WorkloadDefinition
func (cd *ComponentDefinition) SetConditions(c ...runtimev1alpha1.Condition) {
	cd.Status.SetConditions(c...)
}

// GetCondition gets condition from WorkloadDefinition
func (cd *ComponentDefinition) GetCondition(conditionType runtimev1alpha1.ConditionType) runtimev1alpha1.Condition {
	return cd.Status.GetCondition(conditionType)
}

// +kubebuilder:object:root=true

// ComponentDefinitionList contains a list of ComponentDefinition
type ComponentDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ComponentDefinition `json:"items"`
}
