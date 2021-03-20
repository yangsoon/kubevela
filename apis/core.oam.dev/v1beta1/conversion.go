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
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1alpha2"
)

var _ conversion.Convertible = &ComponentDefinition{}
var _ conversion.Convertible = &WorkloadDefinition{}
var _ conversion.Convertible = &Application{}

// ConvertTo converts this ComponentDefinition to the Hub version (v1alpa2).
func (cd *ComponentDefinition) ConvertTo(dstRaw conversion.Hub) error {

	dst := dstRaw.(*v1alpha2.ComponentDefinition)
	klog.Infof("convert %s from %s to %s", cd.Name, cd.APIVersion, dst.APIVersion)
	// set ComponentDefinitionSpec
	dst.SetLabels(cd.Labels)
	dst.SetAnnotations(cd.Annotations)
	dst.Spec.Workload = cd.Spec.Workload
	dst.Spec.ChildResourceKinds = cd.Spec.ChildResourceKinds
	dst.Spec.RevisionLabel = cd.Spec.RevisionLabel
	dst.Spec.PodSpecPath = cd.Spec.PodSpecPath
	dst.Spec.Status = cd.Spec.Status
	dst.Spec.Schematic = cd.Spec.Schematic
	dst.Spec.Extension = cd.Spec.Extension

	// set ComponentDefinitionStatus
	dst.Status.ConditionedStatus = cd.Status.ConditionedStatus
	dst.Status.ConfigMapRef = cd.Status.ConfigMapRef
	return nil
}

// ConvertFrom converts from the Hub version (v1alpa1) to this version (v1beta1).
func (cd *ComponentDefinition) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.ComponentDefinition)
	klog.Infof("convert %s from %s to %s", src.Name, src.APIVersion, cd.APIVersion)
	// set ComponentDefinitionSpec
	cd.SetLabels(src.Labels)
	cd.SetAnnotations(src.Annotations)
	cd.Spec.Workload = src.Spec.Workload
	cd.Spec.ChildResourceKinds = src.Spec.ChildResourceKinds
	cd.Spec.RevisionLabel = src.Spec.RevisionLabel
	cd.Spec.PodSpecPath = src.Spec.PodSpecPath
	cd.Spec.Status = src.Spec.Status
	cd.Spec.Schematic = src.Spec.Schematic
	cd.Spec.Extension = src.Spec.Extension

	return nil
}

// ConvertTo converts this WorkloadDefinition to the Hub version (v1alpa2).
func (wd *WorkloadDefinition) ConvertTo(dstRaw conversion.Hub) error {

	dst := dstRaw.(*v1alpha2.WorkloadDefinition)
	klog.Infof("convert %s from %s to %s", wd.Name, wd.APIVersion, dst.APIVersion)
	// set ComponentDefinitionSpec
	dst.SetLabels(wd.Labels)
	dst.SetAnnotations(wd.Annotations)
	dst.Spec.Reference = wd.Spec.Reference
	dst.Spec.ChildResourceKinds = wd.Spec.ChildResourceKinds
	dst.Spec.RevisionLabel = wd.Spec.RevisionLabel
	dst.Spec.PodSpecPath = wd.Spec.PodSpecPath
	dst.Spec.Status = wd.Spec.Status
	dst.Spec.Schematic = wd.Spec.Schematic
	dst.Spec.Extension = wd.Spec.Extension

	// set ComponentDefinitionStatus
	dst.Status.ConditionedStatus = wd.Status.ConditionedStatus
	return nil
}

// ConvertFrom converts from the Hub version (v1alpa1) to this version (v1beta1).
func (wd *WorkloadDefinition) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.WorkloadDefinition)
	klog.Infof("convert %s from %s to %s", src.Name, src.APIVersion, wd.APIVersion)
	// set ComponentDefinitionSpec
	wd.SetLabels(src.Labels)
	wd.SetAnnotations(src.Annotations)
	wd.Spec.Reference = src.Spec.Reference
	wd.Spec.ChildResourceKinds = src.Spec.ChildResourceKinds
	wd.Spec.RevisionLabel = src.Spec.RevisionLabel
	wd.Spec.PodSpecPath = src.Spec.PodSpecPath
	wd.Spec.Status = src.Spec.Status
	wd.Spec.Schematic = src.Spec.Schematic
	wd.Spec.Extension = src.Spec.Extension

	return nil
}

// ConvertTo converts this WorkloadDefinition to the Hub version (v1alpa2).
func (app *Application) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.Application)
	klog.Infof("convert %s from %s to %s", app.Name, app.APIVersion, dst.APIVersion)
	// set ApplicationSpec
	dst.SetLabels(app.Labels)
	dst.SetAnnotations(app.Annotations)

	if len(app.Spec.Components) > 0 {
		componets := make([]v1alpha2.ApplicationComponent, len(app.Spec.Components))
		for i, component := range app.Spec.Components {
			componets[i] = v1alpha2.ApplicationComponent{
				Name:         component.Name,
				WorkloadType: component.WorkloadType,
				Settings:     component.Properties,
				Traits:       component.Traits,
				Scopes:       component.Scopes,
			}
		}
		dst.Spec.Components = componets
	}
	dst.Spec.RolloutPlan = app.Spec.RolloutPlan

	// set AppStatus
	dst.Status.RollingState = app.Status.RollingState
	dst.Status.Phase = app.Status.Phase
	dst.Status.Components = app.Status.Components
	dst.Status.Services = app.Status.Services
	dst.Status.LatestRevision = app.Status.LatestRevision
	return nil
}

// ConvertFrom converts from the Hub version (v1alpa1) to this version (v1beta1).
func (app *Application) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.Application)
	klog.Infof("convert %s from %s to %s", src.Name, src.APIVersion, app.APIVersion)
	// set ApplicationSpec
	app.SetLabels(app.Labels)
	app.SetAnnotations(app.Annotations)

	if len(src.Spec.Components) > 0 {
		componets := make([]ApplicationComponent, len(src.Spec.Components))
		for i, component := range src.Spec.Components {
			componets[i] = ApplicationComponent{
				Name:         component.Name,
				WorkloadType: component.WorkloadType,
				Properties:   component.Settings,
				Traits:       component.Traits,
				Scopes:       component.Scopes,
			}
		}
		app.Spec.Components = componets
	}
	app.Spec.RolloutPlan = src.Spec.RolloutPlan

	// set AppStatus
	app.Status.RollingState = src.Status.RollingState
	app.Status.Phase = src.Status.Phase
	app.Status.Components = src.Status.Components
	app.Status.Services = src.Status.Services
	app.Status.LatestRevision = src.Status.LatestRevision
	return nil
}
