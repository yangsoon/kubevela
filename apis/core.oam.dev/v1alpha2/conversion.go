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

import (
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this WorkloadDefinition to the Hub version (v1alpa2).
func (app *Application) ConvertTo(dstRaw conversion.Hub) error {

	klog.Infof("ConvertTo Application convert %s from %s to\n", app.Name, app.APIVersion)
	return nil
}

// ConvertFrom converts from the Hub version (v1alpa1) to this version (v1beta1).
func (app *Application) ConvertFrom(srcRaw conversion.Hub) error {

	klog.Infof(" ConvertFrom Application convert %s from %s to\n", app.Name, app.APIVersion)
	return nil
}
