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

package v1alpha2

// Hub marks *v1alpha2.ComponentDefinition as a conversion hub.
func (*ComponentDefinition) Hub() {}

// Hub marks *v1alpha2.WorkloadDefinition as a conversion hub.
func (*WorkloadDefinition) Hub() {}

// Hub marks *v1alpha2.TraitDefinition as a conversion hub.
func (*TraitDefinition) Hub() {}

// Hub marks *v1alpha2.Application as a conversion hub.
func (*Application) Hub() {}

// Hub marks *v1alpha2.AppRollout as a conversion hub.
func (*AppRollout) Hub() {}

// Hub marks *v1alpha2.ApplicationRevision as a conversion hub.
func (*ApplicationRevision) Hub() {}
