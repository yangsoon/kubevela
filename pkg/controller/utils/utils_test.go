package utils

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/stretchr/testify/assert"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	controllerruntime "sigs.k8s.io/controller-runtime"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/common"
	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1alpha2"
	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
	"github.com/oam-dev/kubevela/pkg/oam"
	oamutil "github.com/oam-dev/kubevela/pkg/oam/util"
)

var _ = Describe("utils", func() {
	Context("GetEnabledCapabilities", func() {
		It("disable all", func() {
			disableCaps := "all"
			err := CheckDisabledCapabilities(disableCaps)
			Expect(err).NotTo(HaveOccurred())
		})
		It("disable none", func() {
			disableCaps := ""
			err := CheckDisabledCapabilities(disableCaps)
			Expect(err).NotTo(HaveOccurred())
		})
		It("disable some capabilities", func() {
			disableCaps := "autoscale,route"
			err := CheckDisabledCapabilities(disableCaps)
			Expect(err).NotTo(HaveOccurred())
		})
		It("disable some bad capabilities", func() {
			disableCaps := "abc,def"
			err := CheckDisabledCapabilities(disableCaps)
			Expect(err).To(HaveOccurred())
		})
	})
})

func TestConstructExtract(t *testing.T) {
	tests := []string{"tam1", "test-comp", "xx", "tt-x-x-c"}
	revisionNum := []int{1, 5, 10, 100000}
	for idx, componentName := range tests {
		t.Run(fmt.Sprintf("tests %d for component[%s]", idx, componentName), func(t *testing.T) {
			revisionName := ConstructRevisionName(componentName, int64(revisionNum[idx]))
			got := ExtractComponentName(revisionName)
			if got != componentName {
				t.Errorf("want to get %s from %s but got %s", componentName, revisionName, got)
			}
			revision, _ := ExtractRevision(revisionName)
			if revision != revisionNum[idx] {
				t.Errorf("want to get %d from %s but got %d", revisionNum[idx], revisionName, revision)
			}
		})
	}
	badRevision := []string{"xx", "yy-", "zz-0.1"}
	t.Run(fmt.Sprintf("tests %s for extractRevision", badRevision), func(t *testing.T) {
		for _, revisionName := range badRevision {
			_, err := ExtractRevision(revisionName)
			if err == nil {
				t.Errorf("want to get err from %s but got nil", revisionName)
			}
		}
	})
}

func TestCompareWithRevision(t *testing.T) {
	ctx := context.TODO()
	logger := logging.NewLogrLogger(controllerruntime.Log.WithName("util-test"))
	componentName := "testComp"
	nameSpace := "namespace"
	latestRevision := "revision"
	imageV1 := "wordpress:4.6.1-apache"
	namespaceName := "test"
	cwV1 := v1alpha2.ContainerizedWorkload{
		TypeMeta: v1.TypeMeta{
			Kind:       "ContainerizedWorkload",
			APIVersion: "core.oam.dev/v1alpha2",
		},
		ObjectMeta: v1.ObjectMeta{
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
	baseComp := &v1alpha2.Component{
		TypeMeta: v1.TypeMeta{
			Kind:       "Component",
			APIVersion: "core.oam.dev/v1alpha2",
		}, ObjectMeta: v1.ObjectMeta{
			Name:      "myweb",
			Namespace: namespaceName,
			Labels:    map[string]string{"application.oam.dev": "test"},
		},
		Spec: v1alpha2.ComponentSpec{
			Workload: runtime.RawExtension{
				Object: &cwV1,
			},
		}}

	revisionBase := v12.ControllerRevision{
		ObjectMeta: v1.ObjectMeta{
			Name:      "revisionName",
			Namespace: baseComp.Namespace,
			OwnerReferences: []v1.OwnerReference{
				{
					APIVersion: v1alpha2.SchemeGroupVersion.String(),
					Kind:       v1alpha2.ComponentKind,
					Name:       baseComp.Name,
					UID:        baseComp.UID,
					Controller: pointer.BoolPtr(true),
				},
			},
			Labels: map[string]string{
				"controller.oam.dev/component": baseComp.Name,
			},
		},
		Revision: 2,
		Data:     runtime.RawExtension{Object: baseComp},
	}

	tests := map[string]struct {
		getFunc        test.ObjectFn
		curCompSpec    *v1alpha2.ComponentSpec
		expectedResult bool
		expectedErr    error
	}{
		"compare object": {
			getFunc: func(obj runtime.Object) error {
				o, ok := obj.(*v12.ControllerRevision)
				if !ok {
					t.Errorf("the object %+v is not of type controller revision", o)
				}
				*o = revisionBase
				return nil
			},
			curCompSpec: &v1alpha2.ComponentSpec{
				Workload: runtime.RawExtension{
					Object: baseComp,
				},
			},
			expectedResult: true,
			expectedErr:    nil,
		},
		// TODO: add test cases
		// compare raw with object
		// raw with raw
		// diff in object meta
		// diff in namespace
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tclient := test.MockClient{
				MockGet: test.NewMockGetFn(nil, tt.getFunc),
			}
			same, err := CompareWithRevision(ctx, &tclient, logger, componentName, nameSpace, latestRevision,
				tt.curCompSpec)
			if err != tt.expectedErr {
				t.Errorf("CompareWithRevision() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}
			if same != tt.expectedResult {
				t.Errorf("CompareWithRevision() got = %t, want %t", same, tt.expectedResult)
			}
		})
	}
}

func TestGetAppRevison(t *testing.T) {
	revisionName, latestRevision := GetAppNextRevision(nil)
	assert.Equal(t, revisionName, "")
	assert.Equal(t, latestRevision, int64(0))
	// the first is always 1
	app := &v1beta1.Application{}
	app.Name = "myapp"
	revisionName, latestRevision = GetAppNextRevision(app)
	assert.Equal(t, revisionName, "myapp-v1")
	assert.Equal(t, latestRevision, int64(1))
	app.Status.LatestRevision = &common.Revision{
		Name:     "myapp-v1",
		Revision: 1,
	}
	// we always automatically advance the revision
	revisionName, latestRevision = GetAppNextRevision(app)
	assert.Equal(t, revisionName, "myapp-v2")
	assert.Equal(t, latestRevision, int64(2))
	// we generate new revisions if the app is rolling
	app.SetAnnotations(map[string]string{oam.AnnotationAppRollout: strconv.FormatBool(true)})
	revisionName, latestRevision = GetAppNextRevision(app)
	assert.Equal(t, revisionName, "myapp-v2")
	assert.Equal(t, latestRevision, int64(2))
	app.Status.LatestRevision = &common.Revision{
		Name:     revisionName,
		Revision: latestRevision,
	}
	// try again
	revisionName, latestRevision = GetAppNextRevision(app)
	assert.Equal(t, revisionName, "myapp-v3")
	assert.Equal(t, latestRevision, int64(3))
	app.Status.LatestRevision = &common.Revision{
		Name:     revisionName,
		Revision: latestRevision,
	}
	// remove the annotation and it will still advance
	oamutil.RemoveAnnotations(app, []string{oam.AnnotationAppRollout})
	revisionName, latestRevision = GetAppNextRevision(app)
	assert.Equal(t, revisionName, "myapp-v4")
	assert.Equal(t, latestRevision, int64(4))
}
