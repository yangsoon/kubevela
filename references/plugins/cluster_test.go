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

package plugins

import (
	"context"
	"fmt"
	"io/ioutil"

	"cuelang.org/go/cue"
	"github.com/ghodss/yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/oam-dev/kubevela/pkg/oam/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	corev1beta1 "github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
	"github.com/oam-dev/kubevela/apis/types"
	"github.com/oam-dev/kubevela/pkg/utils/common"
)

const (
	TestDir        = "testdata"
	DeployName     = "deployments.testapps"
	WebserviceName = "webservice.testapps"
)

var _ = Describe("DefinitionFiles", func() {

	deployment := types.Capability{
		Namespace:   "testdef",
		Name:        DeployName,
		Type:        types.TypeComponentDefinition,
		CrdName:     "deployments.apps",
		Description: "description not defined",
		Category:    types.CUECategory,
		Parameters: []types.Parameter{
			{
				Type: cue.ListKind,
				Name: "env",
			},
			{
				Name:     "image",
				Type:     cue.StringKind,
				Default:  "",
				Short:    "i",
				Required: true,
				Usage:    "Which image would you like to use for your service",
			},
			{
				Name:    "port",
				Type:    cue.IntKind,
				Short:   "p",
				Default: int64(8080),
				Usage:   "Which port do you want customer traffic sent to",
			},
		},
		CrdInfo: &types.CRDInfo{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
	}

	websvc := types.Capability{
		Namespace:   "testdef",
		Name:        WebserviceName,
		Type:        types.TypeComponentDefinition,
		Description: "description not defined",
		Category:    types.CUECategory,
		Parameters: []types.Parameter{{
			Name: "env", Type: cue.ListKind,
		}, {
			Name:     "image",
			Type:     cue.StringKind,
			Default:  "",
			Short:    "i",
			Required: true,
			Usage:    "Which image would you like to use for your service",
		}, {
			Name:    "port",
			Type:    cue.IntKind,
			Short:   "p",
			Default: int64(6379),
			Usage:   "Which port do you want customer traffic sent to",
		}},
		CrdName: "deployments.apps",
		CrdInfo: &types.CRDInfo{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
	}

	req, _ := labels.NewRequirement("usecase", selection.Equals, []string{"forplugintest"})
	selector := labels.NewSelector().Add(*req)

	// Notice!!  DefinitionPath Object is Cluster Scope object
	// which means objects created in other DefinitionNamespace will also affect here.
	It("getcomponents", func() {
		workloadDefs, _, err := GetComponentsFromCluster(context.Background(), DefinitionNamespace, common.Args{Config: cfg, Schema: scheme}, selector)
		Expect(err).Should(BeNil())
		logf.Log.Info(fmt.Sprintf("Getting component definitions  %v", workloadDefs))
		for i := range workloadDefs {
			// CueTemplate should always be fulfilled, even those whose CueTemplateURI is assigend,
			By("check CueTemplate is fulfilled")
			Expect(workloadDefs[i].CueTemplate).ShouldNot(BeEmpty())
			workloadDefs[i].CueTemplate = ""
		}
		Expect(cmp.Diff(workloadDefs, []types.Capability{deployment, websvc})).Should(BeEquivalentTo(""))
	})
	It("getall", func() {
		alldef, err := GetCapabilitiesFromCluster(context.Background(), DefinitionNamespace, common.Args{Config: cfg, Schema: scheme}, selector)
		Expect(err).Should(BeNil())
		logf.Log.Info(fmt.Sprintf("Getting all definitions %v", alldef))
		for i := range alldef {
			alldef[i].CueTemplate = ""
		}
		Expect(cmp.Diff(alldef, []types.Capability{deployment, websvc})).Should(BeEquivalentTo(""))
	})
})

var _ = Describe("test GetCapabilityByName", func() {
	var (
		ctx        context.Context
		c          common.Args
		ns         string
		defaultNS  string
		cd1        corev1beta1.ComponentDefinition
		cd2        corev1beta1.ComponentDefinition
		td1        corev1beta1.TraitDefinition
		td2        corev1beta1.TraitDefinition
		component1 string
		component2 string
		trait1     string
		trait2     string
	)
	BeforeEach(func() {
		c = common.Args{
			Client: k8sClient,
			Config: cfg,
			Schema: scheme,
		}
		ctx = context.Background()
		ns = "cluster-test-ns"
		defaultNS = types.DefaultKubeVelaNS
		component1 = "cd1"
		component2 = "cd2"
		trait1 = "td1"
		trait2 = "td2"

		By("create namespace")
		Expect(k8sClient.Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}})).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))
		Expect(k8sClient.Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: defaultNS}})).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

		By("create ComponentDefinition")
		data, _ := ioutil.ReadFile("testdata/componentDef.yaml")
		yaml.Unmarshal(data, &cd1)
		yaml.Unmarshal(data, &cd2)
		cd1.Namespace = ns
		cd1.Name = component1
		Expect(k8sClient.Create(ctx, &cd1)).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

		cd2.Namespace = defaultNS
		cd2.Name = component2
		Expect(k8sClient.Create(ctx, &cd2)).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

		By("create TraitDefinition")
		data, _ = ioutil.ReadFile("testdata/manualscalars.yaml")
		yaml.Unmarshal(data, &td1)
		yaml.Unmarshal(data, &td2)
		td1.Namespace = ns
		td1.Name = trait1
		Expect(k8sClient.Create(ctx, &td1)).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

		td2.Namespace = defaultNS
		td2.Name = trait2
		Expect(k8sClient.Create(ctx, &td2)).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

	})

	It("get capability", func() {
		Context("ComponentDefinition is in the current namespace", func() {
			_, err := GetCapabilityByName(ctx, c, component1, ns)
			Expect(err).Should(BeNil())
		})
		Context("ComponentDefinition is in the default namespace", func() {
			_, err := GetCapabilityByName(ctx, c, component2, ns)
			Expect(err).Should(BeNil())
		})

		Context("TraitDefinition is in the current namespace", func() {
			_, err := GetCapabilityByName(ctx, c, trait1, ns)
			Expect(err).Should(BeNil())
		})
		Context("TraitDefinitionDefinition is in the default namespace", func() {
			_, err := GetCapabilityByName(ctx, c, trait2, ns)
			Expect(err).Should(BeNil())
		})

		Context("capability cloud not be found", func() {
			_, err := GetCapabilityByName(ctx, c, "a-component-definition-not-existed", ns)
			Expect(err).Should(HaveOccurred())
		})
	})
})

var _ = Describe("test GetNamespacedCapabilitiesFromCluster", func() {
	var (
		ctx        context.Context
		c          common.Args
		ns         string
		defaultNS  string
		cd1        corev1beta1.ComponentDefinition
		cd2        corev1beta1.ComponentDefinition
		td1        corev1beta1.TraitDefinition
		td2        corev1beta1.TraitDefinition
		component1 string
		component2 string
		trait1     string
		trait2     string
	)
	BeforeEach(func() {
		c = common.Args{
			Client: k8sClient,
			Config: cfg,
			Schema: scheme,
		}
		ctx = context.Background()
		ns = "cluster-test-ns"
		defaultNS = types.DefaultKubeVelaNS
		component1 = "cd1"
		component2 = "cd2"
		trait1 = "td1"
		trait2 = "td2"

		By("create namespace")
		Expect(k8sClient.Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}})).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))
		Expect(k8sClient.Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: defaultNS}})).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

		By("create ComponentDefinition")
		data, _ := ioutil.ReadFile("testdata/componentDef.yaml")
		yaml.Unmarshal(data, &cd1)
		yaml.Unmarshal(data, &cd2)
		cd1.Namespace = ns
		cd1.Name = component1
		Expect(k8sClient.Create(ctx, &cd1)).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

		cd2.Namespace = defaultNS
		cd2.Name = component2
		Expect(k8sClient.Create(ctx, &cd2)).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

		By("create TraitDefinition")
		data, _ = ioutil.ReadFile("testdata/manualscalars.yaml")
		yaml.Unmarshal(data, &td1)
		yaml.Unmarshal(data, &td2)
		td1.Namespace = ns
		td1.Name = trait1
		Expect(k8sClient.Create(ctx, &td1)).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

		td2.Namespace = defaultNS
		td2.Name = trait2
		Expect(k8sClient.Create(ctx, &td2)).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

	})

	It("get namespaced capabilities", func() {
		Context("found all capabilities", func() {
			capabilities, err := GetNamespacedCapabilitiesFromCluster(ctx, ns, c, nil)
			Expect(len(capabilities)).Should(Equal(4))
			Expect(err).Should(BeNil())
		})

		Context("found two capabilities with a bad namespace", func() {
			capabilities, err := GetNamespacedCapabilitiesFromCluster(ctx, "a-bad-ns", c, nil)
			Expect(len(capabilities)).Should(Equal(2))
			Expect(err).Should(BeNil())
		})

	})
})
