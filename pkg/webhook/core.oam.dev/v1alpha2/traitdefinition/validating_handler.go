package traitdefinition

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1alpha2"
	controller "github.com/oam-dev/kubevela/pkg/controller/core.oam.dev"
	"github.com/oam-dev/kubevela/pkg/oam/discoverymapper"
	"github.com/oam-dev/kubevela/pkg/oam/util"
)

const (
	errValidateDefRef = "error occurs when validating definition reference"

	failInfoDefRefOmitted = "if definition reference is omitted, patch or output with GVK is required"
)

var traitDefGVR = v1alpha2.SchemeGroupVersion.WithResource("traitdefinitions")

// ValidatingHandler handles validation of trait definition
type ValidatingHandler struct {
	Client client.Client
	Mapper discoverymapper.DiscoveryMapper

	// Decoder decodes object
	Decoder *admission.Decoder
	// Validators validate objects
	Validators []TraitDefValidator
}

// TraitDefValidator validate trait definition
type TraitDefValidator interface {
	Validate(context.Context, v1beta1.TraitDefinition) error
}

// TraitDefValidatorFn implements TraitDefValidator
type TraitDefValidatorFn func(context.Context, v1beta1.TraitDefinition) error

// Validate implements TraitDefValidator method
func (fn TraitDefValidatorFn) Validate(ctx context.Context, td v1beta1.TraitDefinition) error {
	return fn(ctx, td)
}

var _ admission.Handler = &ValidatingHandler{}

// Handle validate trait definition
func (h *ValidatingHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	obj := &v1beta1.TraitDefinition{}
	if req.Resource.String() != traitDefGVR.String() {
		return admission.Errored(http.StatusBadRequest, fmt.Errorf("expect resource to be %s", traitDefGVR))
	}

	if req.Operation == admissionv1beta1.Create || req.Operation == admissionv1beta1.Update {
		err := h.Decoder.Decode(req, obj)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
		klog.Info("validating ", " name: ", obj.Name, " operation: ", string(req.Operation))
		for _, validator := range h.Validators {
			if err := validator.Validate(ctx, *obj); err != nil {
				klog.Info("validation failed ", " name: ", obj.Name, " errMsgi: ", err.Error())
				return admission.Denied(err.Error())
			}
		}
		klog.Info("validation passed ", " name: ", obj.Name, " operation: ", string(req.Operation))
	}
	return admission.ValidationResponse(true, "")
}

var _ inject.Client = &ValidatingHandler{}

// InjectClient injects the client into the ValidatingHandler
func (h *ValidatingHandler) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}

var _ admission.DecoderInjector = &ValidatingHandler{}

// InjectDecoder injects the decoder into the ValidatingHandler
func (h *ValidatingHandler) InjectDecoder(d *admission.Decoder) error {
	h.Decoder = d
	return nil
}

// RegisterValidatingHandler will register TraitDefinition validation to webhook
func RegisterValidatingHandler(mgr manager.Manager, args controller.Args) {
	server := mgr.GetWebhookServer()
	server.Register("/validating-core-oam-dev-v1alpha2-traitdefinitions", &webhook.Admission{Handler: &ValidatingHandler{
		Mapper: args.DiscoveryMapper,
		Validators: []TraitDefValidator{
			TraitDefValidatorFn(ValidateDefinitionReference),
			// add more validators here
		},
	}})
}

// ValidateDefinitionReference validates whether the trait definition is valid if
// its `.spec.reference` field is unset.
// It's valid if
// it has at least one output, and all outputs must have GVK
// or it has no output but has a patch
// or it has a patch and outputs, and all outputs must have GVK
// TODO(roywang) currently we only validate whether it contains CUE template.
// Further validation, e.g., output with GVK, valid patch, etc, remains to be done.
func ValidateDefinitionReference(_ context.Context, td v1beta1.TraitDefinition) error {
	if len(td.Spec.Reference.Name) > 0 {
		return nil
	}
	tmp, err := util.NewTemplate(td.Spec.Schematic, td.Spec.Status, td.Spec.Extension)
	if err != nil {
		return errors.Wrap(err, errValidateDefRef)
	}
	if len(tmp.TemplateStr) == 0 {
		return errors.New(failInfoDefRefOmitted)
	}
	return nil
}
