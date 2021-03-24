package application

import (
	"context"
	"math/rand"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/common"
	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1alpha2"
	"github.com/oam-dev/kubevela/pkg/controller/utils"
)

var _ = Describe("Test Application apply", func() {
	var handler appHandler
	var app *v1alpha2.Application
	var namespaceName string
	var ns corev1.Namespace

	BeforeEach(func() {
		ctx := context.TODO()
		namespaceName = "apply-test-" + strconv.Itoa(rand.Intn(1000))
		ns = corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespaceName,
			},
		}
		app = &v1alpha2.Application{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Application",
				APIVersion: "core.oam.dev/v1alpha2",
			},
		}
		app.Namespace = namespaceName
		app.Spec = v1alpha2.ApplicationSpec{
			Components: []v1alpha2.ApplicationComponent{{
				WorkloadType: "webservice",
				Name:         "express-server",
				Scopes:       map[string]string{"healthscopes.core.oam.dev": "myapp-default-health"},
				Settings: runtime.RawExtension{
					Raw: []byte(`{"image": "oamdev/testapp:v1", "cmd": ["node", "server.js"]}`),
				},
				Traits: []common.ApplicationTrait{{
					Name: "route",
					Properties: runtime.RawExtension{
						Raw: []byte(`{"domain": "example.com", "http":{"/": 8080}}`),
					},
				},
				},
			}},
		}
		handler = appHandler{
			r:      reconciler,
			app:    app,
			logger: reconciler.Log.WithValues("application", "unit-test"),
		}

		By("Create the Namespace for test")
		Expect(k8sClient.Create(ctx, &ns)).Should(Succeed())
	})

	AfterEach(func() {
		By("[TEST] Clean up resources after an integration test")
		Expect(k8sClient.Delete(context.TODO(), &ns)).Should(Succeed())
	})

	It("Test update or create component", func() {
		ctx := context.TODO()
		By("[TEST] Setting up the testing environment")
		imageV1 := "wordpress:4.6.1-apache"
		imageV2 := "wordpress:4.6.2-apache"
		cwV1 := v1alpha2.ContainerizedWorkload{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ContainerizedWorkload",
				APIVersion: "core.oam.dev/v1alpha2",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespaceName,
			},
			Spec: v1alpha2.ContainerizedWorkloadSpec{
				Containers: []v1alpha2.Container{
					{
						Name:  "wordpress",
						Image: imageV1,
						Ports: []v1alpha2.ContainerPort{
							{
								Name: "wordpress",
								Port: 80,
							},
						},
					},
				},
			},
		}
		component := &v1alpha2.Component{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Component",
				APIVersion: "core.oam.dev/v1alpha2",
			}, ObjectMeta: metav1.ObjectMeta{
				Name:      "myweb",
				Namespace: namespaceName,
				Labels:    map[string]string{"application.oam.dev": "test"},
			},
			Spec: v1alpha2.ComponentSpec{
				Workload: runtime.RawExtension{
					Object: &cwV1,
				},
			}}

		By("[TEST] Creating a component the first time")
		// take a copy so the component's workload still uses object instead of raw data
		// just like the way we use it in prod. The raw data will be filled by the k8s for some reason.
		revision, err := handler.createOrUpdateComponent(ctx, component.DeepCopy())
		By("verify that the revision is the set correctly and newRevision is true")
		Expect(err).ShouldNot(HaveOccurred())
		// verify the revision actually contains the right component
		Expect(utils.CompareWithRevision(ctx, handler.r, logging.NewLogrLogger(handler.logger), component.GetName(),
			component.GetNamespace(), revision, &component.Spec)).Should(BeTrue())
		preRevision := revision

		By("[TEST] update the component without any changes (mimic reconcile behavior)")
		revision, err = handler.createOrUpdateComponent(ctx, component.DeepCopy())
		By("verify that the revision is the same and newRevision is false")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(revision).Should(BeIdenticalTo(preRevision))

		By("[TEST] update the component")
		// modify the component spec through object
		cwV2 := cwV1.DeepCopy()
		cwV2.Spec.Containers[0].Image = imageV2
		component.Spec.Workload.Object = cwV2
		revision, err = handler.createOrUpdateComponent(ctx, component.DeepCopy())
		By("verify that the revision is changed and newRevision is true")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(revision).ShouldNot(BeIdenticalTo(preRevision))
		Expect(utils.CompareWithRevision(ctx, handler.r, logging.NewLogrLogger(handler.logger), component.GetName(),
			component.GetNamespace(), revision, &component.Spec)).Should(BeTrue())
		// revision increased
		Expect(strings.Compare(revision, preRevision) > 0).Should(BeTrue())
	})

})
