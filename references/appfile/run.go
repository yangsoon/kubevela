package appfile

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1alpha2"
	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
	"github.com/oam-dev/kubevela/pkg/oam"
)

// Run will deploy OAM objects and other assistant K8s Objects including ConfigMap, OAM Scope Custom Resource.
func Run(ctx context.Context, client client.Client, app *v1beta1.Application, assistantObjects []oam.Object) error {
	if err := CreateOrUpdateObjects(ctx, client, assistantObjects); err != nil {
		return err
	}
	if app != nil {
		if err := CreateOrUpdateApplication(ctx, client, app); err != nil {
			return err
		}
	}
	return nil
}

// CreateOrUpdateObjects will create all scopes
func CreateOrUpdateObjects(ctx context.Context, client client.Client, objects []oam.Object) error {
	for _, obj := range objects {
		key := ctypes.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}
		u := &unstructured.Unstructured{}
		u.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
		err := client.Get(ctx, key, u)
		if err == nil {
			obj.SetResourceVersion(u.GetResourceVersion())
			fmt.Println("Updating: ", u.GetObjectKind().GroupVersionKind().String(), "in", u.GetNamespace())
			if err = client.Update(ctx, obj); err != nil {
				return err
			}
			continue
		}
		if !apierrors.IsNotFound(err) {
			return err
		}
		fmt.Println("Creating: ", u.GetObjectKind().GroupVersionKind().String(), "in", u.GetNamespace())
		if err = client.Create(ctx, obj); err != nil {
			return err
		}
	}
	return nil
}

// CreateOrUpdateApplication will create if not exist and update if exists.
func CreateOrUpdateApplication(ctx context.Context, client client.Client, app *v1beta1.Application) error {
	var geta v1alpha2.Application
	key := ctypes.NamespacedName{Name: app.Name, Namespace: app.Namespace}
	var exist = true
	if err := client.Get(ctx, key, &geta); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		exist = false
	}
	if !exist {
		return client.Create(ctx, app)
	}
	app.ResourceVersion = geta.ResourceVersion
	return client.Update(ctx, app)
}
